package grpcserver_test

import (
	"context"
	"testing"

	"go-market/internal/platform/grpcserver"
	"go-market/internal/platform/requestid"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func TestRequestIDUnaryServerInterceptor(t *testing.T) {
	const validRequestID = "eea47e6f-9c35-4cc3-83fc-3d576f900749"

	tests := []struct {
		name          string
		incomingID    string
		wantPreserved bool
	}{
		{
			name:          "valid request ID",
			incomingID:    validRequestID,
			wantPreserved: true,
		},
		{
			name:       "missing request ID",
			incomingID: "",
		},
		{
			name:       "invalid request ID",
			incomingID: "not-a-uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			if tt.incomingID != "" {
				ctx = metadata.NewIncomingContext(
					ctx,
					metadata.Pairs(
						requestid.GRPCMetadataKey,
						tt.incomingID,
					),
				)
			}

			var handlerRequestID string
			handlerCalled := false

			_, err := grpcserver.RequestIDUnaryServerInterceptor(
				ctx,
				nil,
				nil,
				func(ctx context.Context, _ any) (any, error) {
					handlerCalled = true

					id, ok := requestid.FromContext(ctx)
					if !ok {
						t.Error("request ID is missing from handler context")
					}

					handlerRequestID = id
					return nil, nil
				},
			)
			if err != nil {
				t.Fatalf("interceptor returned an error: %v", err)
			}

			if !handlerCalled {
				t.Fatal("handler was not called")
			}

			if _, err := uuid.Parse(handlerRequestID); err != nil {
				t.Errorf(
					"handler request ID %q is not a valid UUID",
					handlerRequestID,
				)
			}

			if tt.wantPreserved && handlerRequestID != tt.incomingID {
				t.Errorf(
					"handler request ID = %q, want %q",
					handlerRequestID,
					tt.incomingID,
				)
			}

			if !tt.wantPreserved &&
				tt.incomingID != "" &&
				handlerRequestID == tt.incomingID {
				t.Errorf(
					"invalid request ID %q was not replaced",
					tt.incomingID,
				)
			}
		})
	}
}
