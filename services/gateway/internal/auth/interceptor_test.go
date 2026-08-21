package auth_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	catalogv1 "go-market/gen/go/catalog/v1"
	"go-market/internal/platform/identity"
	gatewayauth "go-market/services/gateway/internal/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type verifierStub struct {
	principal identity.Principal
	err       error
	calls     int
	rawToken  string
}

func (v *verifierStub) Verify(rawToken string) (identity.Principal, error) {
	v.calls++
	v.rawToken = rawToken

	return v.principal, v.err
}

func TestInterceptorAllowsPublicMethod(t *testing.T) {
	verifier := &verifierStub{}
	invoked := false

	interceptor := gatewayauth.NewInterceptor(
		testLogger(),
		verifier,
	)

	err := interceptor.UnaryClient(
		context.Background(),
		catalogv1.CatalogService_ListProducts_FullMethodName,
		nil,
		nil,
		nil,
		func(
			context.Context,
			string,
			any,
			any,
			*grpc.ClientConn,
			...grpc.CallOption,
		) error {
			invoked = true
			return nil
		},
	)

	if err != nil {
		t.Fatalf("UnaryClient() error = %v", err)
	}

	if !invoked {
		t.Fatal("public RPC was not invoked")
	}

	if verifier.calls != 0 {
		t.Fatal("verifier was called without a token")
	}
}

func TestInterceptorRejectsPrivateMethod(t *testing.T) {
	verifier := &verifierStub{}
	invoked := false

	interceptor := gatewayauth.NewInterceptor(
		testLogger(),
		verifier,
	)

	err := interceptor.UnaryClient(
		context.Background(),
		healthv1.Health_List_FullMethodName,
		nil,
		nil,
		nil,
		func(
			context.Context,
			string,
			any,
			any,
			*grpc.ClientConn,
			...grpc.CallOption,
		) error {
			invoked = true
			return nil
		},
	)

	if got, want := status.Code(err), codes.Unauthenticated; got != want {
		t.Fatalf("gRPC code = %s, want %s", got, want)
	}

	if invoked {
		t.Fatal("private RPC was invoked")
	}
}

func TestInterceptorAllowsAuthenticatedPrivateMethod(t *testing.T) {
	verifier := &verifierStub{
		principal: identity.Principal{
			UserID: "user-123",
			Roles:  []string{"customer"},
		},
	}

	interceptor := gatewayauth.NewInterceptor(
		testLogger(),
		verifier,
	)

	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(
			"authorization", "Bearer signed-token",
			"x-user-id", "spoofed-user",
		),
	)

	var userIDs []string
	var roles []string

	err := interceptor.UnaryClient(
		ctx,
		healthv1.Health_List_FullMethodName,
		nil,
		nil,
		nil,
		func(
			ctx context.Context,
			_ string,
			_ any,
			_ any,
			_ *grpc.ClientConn,
			_ ...grpc.CallOption,
		) error {
			md, _ := metadata.FromOutgoingContext(ctx)

			userIDs = md.Get("x-user-id")
			roles = md.Get("x-user-roles")

			return nil
		},
	)

	if err != nil {
		t.Fatalf("UnaryClient() error = %v", err)
	}

	if verifier.rawToken != "signed-token" {
		t.Errorf(
			"verified token = %q, want %q",
			verifier.rawToken,
			"signed-token",
		)
	}

	if len(userIDs) != 1 || userIDs[0] != "user-123" {
		t.Errorf("user metadata = %v", userIDs)
	}

	if len(roles) != 1 || roles[0] != "customer" {
		t.Errorf("roles metadata = %v", roles)
	}
}

func testLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)
}
