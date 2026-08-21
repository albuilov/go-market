package grpcmiddleware_test

import (
	"context"
	"testing"
	"time"

	"go-market/pkg/grpcmiddleware"

	"google.golang.org/grpc"
)

func TestTimeoutUnaryClientInterceptorAddsDeadline(t *testing.T) {
	const timeout = time.Second

	interceptor := grpcmiddleware.TimeoutUnaryClientInterceptor(timeout)

	err := interceptor(
		context.Background(),
		"/catalog.v1.CatalogService/ListProducts",
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
			deadline, ok := ctx.Deadline()
			if !ok {
				t.Fatal("invocation context does not have a deadline")
			}

			remaining := time.Until(deadline)
			if remaining <= 0 || remaining > timeout {
				t.Errorf(
					"deadline remaining time = %s, want a value in (0, %s]",
					remaining,
					timeout,
				)
			}

			return nil
		},
	)
	if err != nil {
		t.Fatalf("interceptor returned an error: %v", err)
	}
}

func TestTimeoutUnaryClientInterceptorPreservesEarlierDeadline(t *testing.T) {
	parentContext, cancel := context.WithTimeout(
		context.Background(),
		100*time.Millisecond,
	)
	defer cancel()

	parentDeadline, ok := parentContext.Deadline()
	if !ok {
		t.Fatal("parent context does not have a deadline")
	}

	interceptor := grpcmiddleware.TimeoutUnaryClientInterceptor(
		time.Second,
	)

	err := interceptor(
		parentContext,
		"/catalog.v1.CatalogService/ListProducts",
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
			deadline, ok := ctx.Deadline()
			if !ok {
				t.Fatal("invocation context does not have a deadline")
			}

			if !deadline.Equal(parentDeadline) {
				t.Errorf(
					"invocation deadline = %s, want parent deadline %s",
					deadline,
					parentDeadline,
				)
			}

			return nil
		},
	)
	if err != nil {
		t.Fatalf("interceptor returned an error: %v", err)
	}
}
