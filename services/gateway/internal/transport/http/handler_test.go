package http

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	catalogv1 "go-market/gen/go/catalog/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func newTestLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)
}

type catalogClientStub struct {
	response *catalogv1.ListProductsResponse
	err      error
	onCall   func(context.Context)
}

func (s catalogClientStub) ListProducts(
	ctx context.Context,
	_ *catalogv1.ListProductsRequest,
	_ ...grpc.CallOption,
) (*catalogv1.ListProductsResponse, error) {
	if s.onCall != nil {
		s.onCall(ctx)
	}

	return s.response, s.err
}

func TestHandlerListProducts(t *testing.T) {
	handler, err := NewHandler(
		context.Background(),
		newTestLogger(),
		catalogClientStub{
			response: &catalogv1.ListProductsResponse{
				Products: []*catalogv1.Product{
					{
						Id:              "product-1",
						Name:            "Keyboard",
						PriceMinorUnits: 129900,
						CurrencyCode:    "RUB",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if !strings.Contains(
		recorder.Body.String(),
		`"price_minor_units":"129900"`,
	) {
		t.Errorf("response has unexpected JSON format: %s", recorder.Body.String())
	}

	if got, want := recorder.Code, http.StatusOK; got != want {
		t.Fatalf("status code = %d, want %d", got, want)
	}

	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	if !strings.Contains(recorder.Body.String(), `"id":"product-1"`) {
		t.Errorf("response body does not contain product: %s", recorder.Body.String())
	}
}

func TestHandlerMapsGRPCErrorToHTTP(t *testing.T) {
	handler, err := NewHandler(
		context.Background(),
		newTestLogger(),
		catalogClientStub{
			err: status.Error(codes.NotFound, "products not found"),
		},
	)
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/products",
		nil,
	)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if got, want := recorder.Code, http.StatusNotFound; got != want {
		t.Fatalf("status code = %d, want %d", got, want)
	}

	if !strings.Contains(recorder.Body.String(), `"code":"NotFound"`) {
		t.Errorf("unexpected response body: %s", recorder.Body.String())
	}
}

func TestHandlerForwardsRequestIDToGRPC(t *testing.T) {
	const requestID = "eea47e6f-9c35-4cc3-83fc-3d576f900749"

	var grpcRequestID string

	handler, err := NewHandler(
		context.Background(),
		newTestLogger(),
		catalogClientStub{
			response: &catalogv1.ListProductsResponse{},
			onCall: func(ctx context.Context) {
				md, ok := metadata.FromOutgoingContext(ctx)
				if !ok {
					return
				}

				values := md.Get("x-request-id")
				if len(values) > 0 {
					grpcRequestID = values[0]
				}
			},
		},
	)
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/products",
		nil,
	)
	request.Header.Set("X-Request-ID", requestID)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	if got := recorder.Header().Get("X-Request-ID"); got != requestID {
		t.Errorf("response request ID = %q, want %q", got, requestID)
	}

	if grpcRequestID != requestID {
		t.Errorf("gRPC request ID = %q, want %q", grpcRequestID, requestID)
	}
}

func TestHandlerDoesNotExposeInternalGRPCError(t *testing.T) {
	handler, err := NewHandler(
		context.Background(),
		newTestLogger(),
		catalogClientStub{
			err: status.Error(
				codes.Internal,
				"database connection contains sensitive details",
			),
		},
	)
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/products",
		nil,
	)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if got, want := recorder.Code, http.StatusInternalServerError; got != want {
		t.Errorf("status code = %d, want %d", got, want)
	}

	body := recorder.Body.String()

	if !strings.Contains(body, `"code":"Internal"`) {
		t.Errorf("response does not contain error code: %s", body)
	}

	if !strings.Contains(body, `"message":"internal server error"`) {
		t.Errorf("response does not contain public error message: %s", body)
	}

	if strings.Contains(body, "sensitive details") {
		t.Errorf("internal error details were exposed: %s", body)
	}
}
