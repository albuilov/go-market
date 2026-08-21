package app

import (
	"context"
	"testing"

	catalogv1 "go-market/gen/go/catalog/v1"

	"buf.build/go/protovalidate"
	grpcprotovalidate "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestValidationInterceptor(t *testing.T) {
	validator, err := protovalidate.New()
	if err != nil {
		t.Fatalf("create validator: %v", err)
	}

	interceptor := grpcprotovalidate.UnaryServerInterceptor(validator)

	tests := []struct {
		name              string
		request           *catalogv1.ListProductsRequest
		wantCode          codes.Code
		wantHandlerCalled bool
	}{
		{
			name:              "page size omitted",
			request:           &catalogv1.ListProductsRequest{},
			wantCode:          codes.OK,
			wantHandlerCalled: true,
		},
		{
			name: "valid page size",
			request: &catalogv1.ListProductsRequest{
				PageSize: proto.Int32(20),
			},
			wantCode:          codes.OK,
			wantHandlerCalled: true,
		},
		{
			name: "page size is zero",
			request: &catalogv1.ListProductsRequest{
				PageSize: proto.Int32(0),
			},
			wantCode:          codes.InvalidArgument,
			wantHandlerCalled: false,
		},
		{
			name: "page size exceeds limit",
			request: &catalogv1.ListProductsRequest{
				PageSize: proto.Int32(101),
			},
			wantCode:          codes.InvalidArgument,
			wantHandlerCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled := false

			_, err := interceptor(
				context.Background(),
				tt.request,
				&grpc.UnaryServerInfo{
					FullMethod: catalogv1.CatalogService_ListProducts_FullMethodName,
				},
				func(context.Context, any) (any, error) {
					handlerCalled = true
					return &catalogv1.ListProductsResponse{}, nil
				},
			)

			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("status code = %s, want %s", got, tt.wantCode)
			}

			if handlerCalled != tt.wantHandlerCalled {
				t.Errorf(
					"handler called = %t, want %t",
					handlerCalled,
					tt.wantHandlerCalled,
				)
			}
		})
	}
}
