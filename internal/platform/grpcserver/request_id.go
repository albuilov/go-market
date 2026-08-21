package grpcserver

import (
	"context"

	"go-market/internal/platform/requestid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestIDUnaryServerInterceptor переносит идентификатор запроса
// из входящей gRPC metadata в context.
// Если идентификатор отсутствует или некорректен, создаётся новый UUID.
func RequestIDUnaryServerInterceptor(
	ctx context.Context,
	request any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	values := metadata.ValueFromIncomingContext(
		ctx,
		requestid.GRPCMetadataKey,
	)

	var id string
	if len(values) > 0 {
		id = values[0]
	}

	id = requestid.Normalize(id)
	ctx = requestid.WithID(ctx, id)

	return handler(ctx, request)
}
