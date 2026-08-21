package auth

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"go-market/internal/platform/identity"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TokenVerifier interface {
	Verify(rawToken string) (identity.Principal, error)
}

type Interceptor struct {
	logger   *slog.Logger
	policy   Policy
	verifier TokenVerifier
}

func NewInterceptor(
	logger *slog.Logger,
	verifier TokenVerifier,
) *Interceptor {
	return &Interceptor{
		logger:   logger,
		policy:   Policy{},
		verifier: verifier,
	}
}

func (i *Interceptor) UnaryClient(
	ctx context.Context,
	method string,
	request any,
	reply any,
	connection *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	options ...grpc.CallOption,
) error {
	public, err := i.policy.IsPublic(method)
	if err != nil {
		i.logger.ErrorContext(
			ctx,
			"failed to resolve RPC authorization policy",
			slog.String("method", method),
			slog.Any("error", err),
		)

		return status.Error(codes.Internal, "internal server error")
	}

	ctx = identity.RemoveFromOutgoing(ctx)

	rawToken, tokenProvided, err := bearerToken(ctx)
	if err != nil {
		return status.Error(
			codes.Unauthenticated,
			"invalid authorization header",
		)
	}

	if !tokenProvided {
		if public {
			return invoker(
				ctx,
				method,
				request,
				reply,
				connection,
				options...,
			)
		}

		return status.Error(
			codes.Unauthenticated,
			"authentication required",
		)
	}

	principal, err := i.verifier.Verify(rawToken)
	if err != nil {
		return status.Error(
			codes.Unauthenticated,
			"invalid access token",
		)
	}

	ctx = identity.AddToOutgoing(ctx, principal)

	return invoker(
		ctx,
		method,
		request,
		reply,
		connection,
		options...,
	)
}

func bearerToken(ctx context.Context) (string, bool, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", false, nil
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return "", false, nil
	}

	if len(values) != 1 {
		return "", true, errors.New(
			"multiple authorization headers",
		)
	}

	parts := strings.Fields(values[0])
	if len(parts) != 2 ||
		!strings.EqualFold(parts[0], "Bearer") ||
		parts[1] == "" {
		return "", true, errors.New(
			"authorization header must use Bearer scheme",
		)
	}

	return parts[1], true, nil
}
