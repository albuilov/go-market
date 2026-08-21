package http

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	catalogv1 "go-market/gen/go/catalog/v1"
	platformhealth "go-market/internal/platform/health"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
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

type healthClientStub struct {
	response *healthv1.HealthCheckResponse
	err      error
}

func testReadiness(client healthClientStub) platformhealth.Checker {
	return platformhealth.NewGRPCChecker(
		client,
		catalogv1.CatalogService_ServiceDesc.ServiceName,
	)
}

func (s healthClientStub) Check(
	context.Context,
	*healthv1.HealthCheckRequest,
	...grpc.CallOption,
) (*healthv1.HealthCheckResponse, error) {
	if s.response == nil && s.err == nil {
		return &healthv1.HealthCheckResponse{
			Status: healthv1.HealthCheckResponse_SERVING,
		}, nil
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
						Id:           "product-1",
						Name:         "Keyboard",
						Price:        1299,
						CurrencyCode: "RUB",
					},
				},
			},
		},
		testReadiness(healthClientStub{}),
	)
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if !strings.Contains(
		recorder.Body.String(),
		`"price":1299`,
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
		testReadiness(healthClientStub{}),
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

func TestHandlerLogsGRPCErrorOnce(t *testing.T) {
	var output bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&output, nil))

	handler, err := NewHandler(
		context.Background(),
		logger,
		catalogClientStub{
			err: status.Error(codes.NotFound, "products not found"),
		},
		testReadiness(healthClientStub{}),
	)
	if err != nil {
		t.Fatalf("NewHandler() error = %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	logs := output.String()
	if got := strings.Count(logs, `"msg":"HTTP request completed"`); got != 1 {
		t.Errorf("HTTP completion log count = %d, want 1; logs: %s", got, logs)
	}
	if strings.Contains(logs, "gateway request failed") {
		t.Errorf("duplicate gateway error log found: %s", logs)
	}
	if !strings.Contains(logs, `"level":"WARN"`) ||
		!strings.Contains(logs, `"status":404`) ||
		!strings.Contains(logs, "products not found") {
		t.Errorf("completion log does not contain error context: %s", logs)
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
		testReadiness(healthClientStub{}),
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
		testReadiness(healthClientStub{}),
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
