package grpcmiddleware_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"go-market/pkg/grpcmiddleware"
	"go-market/pkg/requestid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRecoveryUnaryServerInterceptor(t *testing.T) {
	const (
		requestID  = "eea47e6f-9c35-4cc3-83fc-3d576f900749"
		fullMethod = "/catalog.v1.CatalogService/ListProducts"
	)

	var logs bytes.Buffer

	logger := slog.New(
		slog.NewJSONHandler(&logs, nil),
	)

	ctx := requestid.WithID(
		context.Background(),
		requestID,
	)

	interceptor := grpcmiddleware.RecoveryUnaryServerInterceptor(logger)

	response, err := interceptor(
		ctx,
		nil,
		&grpc.UnaryServerInfo{
			FullMethod: fullMethod,
		},
		func(context.Context, any) (any, error) {
			panic("unexpected failure")
		},
	)

	if response != nil {
		t.Errorf("response = %#v, want nil", response)
	}

	if got, want := status.Code(err), codes.Internal; got != want {
		t.Errorf("status code = %s, want %s", got, want)
	}

	if got, want := status.Convert(err).Message(), "internal server error"; got != want {
		t.Errorf("status message = %q, want %q", got, want)
	}

	logOutput := logs.String()

	for _, value := range []string{
		requestID,
		fullMethod,
		"unexpected failure",
		"stack_trace",
	} {
		if !strings.Contains(logOutput, value) {
			t.Errorf("log does not contain %q: %s", value, logOutput)
		}
	}

	if strings.Contains(status.Convert(err).Message(), "unexpected failure") {
		t.Error("internal panic details were exposed to the client")
	}
}
