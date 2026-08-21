package grpcmiddleware

import (
	"context"
	"log/slog"
	"runtime/debug"

	"go-market/pkg/requestid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryUnaryServerInterceptor перехватывает panic во время обработки
// unary gRPC-вызова и преобразует его в безопасную ошибку Internal.
func RecoveryUnaryServerInterceptor(
	logger *slog.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		request any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (response any, err error) {
		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}

			id, _ := requestid.FromContext(ctx)

			logger.ErrorContext(
				ctx,
				"panic recovered in gRPC handler",
				slog.String("request_id", id),
				slog.String("method", info.FullMethod),
				slog.Any("panic", recovered),
				slog.String("stack_trace", string(debug.Stack())),
			)

			response = nil
			err = status.Error(
				codes.Internal,
				"internal server error",
			)
		}()

		return handler(ctx, request)
	}
}
