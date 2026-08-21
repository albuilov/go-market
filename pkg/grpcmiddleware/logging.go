package grpcmiddleware

import (
	"context"
	"log/slog"
	"time"

	"go-market/pkg/requestid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingUnaryServerInterceptor записывает результат выполнения unary gRPC-вызова.
func LoggingUnaryServerInterceptor(
	logger *slog.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		request any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startedAt := time.Now()

		response, err := handler(ctx, request)
		code := status.Code(err)

		id, _ := requestid.FromContext(ctx)

		logger.Log(
			ctx,
			grpcLogLevel(code),
			"gRPC request completed",
			slog.String("request_id", id),
			slog.String("method", info.FullMethod),
			slog.String("grpc_code", code.String()),
			slog.Duration("duration", time.Since(startedAt)),
		)

		return response, err
	}
}

func grpcLogLevel(code codes.Code) slog.Level {
	switch code {
	case codes.OK, codes.Canceled:
		return slog.LevelInfo

	case codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.Unauthenticated,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.Aborted,
		codes.OutOfRange,
		codes.DeadlineExceeded:
		return slog.LevelWarn

	default:
		return slog.LevelError
	}
}
