package grpcclient

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TimeoutUnaryClientInterceptor ограничивает длительность unary gRPC-вызова.
// Если исходный context имеет более короткий deadline, он сохраняется.
func TimeoutUnaryClientInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		request any,
		reply any,
		connection *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		options ...grpc.CallOption,
	) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return invoker(
			ctx,
			method,
			request,
			reply,
			connection,
			options...,
		)
	}
}
