package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func TestHealthEndpoints(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		health     healthClientStub
		wantStatus int
		wantBody   string
	}{
		{
			name: "liveness does not depend on catalog",
			path: "/healthz",
			health: healthClientStub{
				err: status.Error(
					codes.Unavailable,
					"catalog is unavailable",
				),
			},
			wantStatus: http.StatusOK,
			wantBody:   `"status":"ok"`,
		},
		{
			name:       "readiness succeeds when catalog is serving",
			path:       "/readyz",
			health:     healthClientStub{},
			wantStatus: http.StatusOK,
			wantBody:   `"status":"ok"`,
		},
		{
			name: "readiness fails when catalog is not serving",
			path: "/readyz",
			health: healthClientStub{
				response: &healthv1.HealthCheckResponse{
					Status: healthv1.HealthCheckResponse_NOT_SERVING,
				},
			},
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   `"status":"not_ready"`,
		},
		{
			name: "readiness fails on catalog error",
			path: "/readyz",
			health: healthClientStub{
				err: status.Error(
					codes.Unavailable,
					"catalog is unavailable",
				),
			},
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   `"status":"not_ready"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := NewHandler(
				context.Background(),
				newTestLogger(),
				catalogClientStub{},
				testReadiness(tt.health),
			)
			if err != nil {
				t.Fatalf("NewHandler() error = %v", err)
			}

			request := httptest.NewRequest(
				http.MethodGet,
				tt.path,
				nil,
			)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			if got := recorder.Code; got != tt.wantStatus {
				t.Errorf(
					"status code = %d, want %d",
					got,
					tt.wantStatus,
				)
			}

			if body := recorder.Body.String(); !strings.Contains(
				body,
				tt.wantBody,
			) {
				t.Errorf(
					"response body %q does not contain %q",
					body,
					tt.wantBody,
				)
			}

			if got := recorder.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf(
					"Content-Type = %q, want application/json",
					got,
				)
			}

			if got := recorder.Header().Get("Cache-Control"); got != "no-store" {
				t.Errorf(
					"Cache-Control = %q, want no-store",
					got,
				)
			}

			if got := recorder.Header().Get("X-Request-ID"); got == "" {
				t.Error("X-Request-ID response header is missing")
			}
		})
	}
}
