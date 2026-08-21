package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"go-market/pkg/requestid"

	"google.golang.org/grpc"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

const readinessCheckTimeout = time.Second

type grpcHealthClient interface {
	Check(
		ctx context.Context,
		request *healthv1.HealthCheckRequest,
		options ...grpc.CallOption,
	) (*healthv1.HealthCheckResponse, error)
}

func livenessHandler(
	w http.ResponseWriter,
	_ *http.Request,
) {
	writeHealthResponse(
		w,
		http.StatusOK,
		`{"status":"ok"}`,
	)
}

func readinessHandler(
	logger *slog.Logger,
	client grpcHealthClient,
	service string,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(
			r.Context(),
			readinessCheckTimeout,
		)
		defer cancel()

		response, err := client.Check(
			ctx,
			&healthv1.HealthCheckRequest{
				Service: service,
			},
		)

		statusValue := healthv1.HealthCheckResponse_UNKNOWN
		if response != nil {
			statusValue = response.GetStatus()
		}

		if err != nil ||
			statusValue != healthv1.HealthCheckResponse_SERVING {
			id, _ := requestid.FromContext(r.Context())

			logger.WarnContext(
				r.Context(),
				"gateway readiness check failed",
				slog.String("request_id", id),
				slog.String("service", service),
				slog.String("status", statusValue.String()),
				slog.Any("error", err),
			)

			writeHealthResponse(
				w,
				http.StatusServiceUnavailable,
				`{"status":"not_ready"}`,
			)
			return
		}

		writeHealthResponse(
			w,
			http.StatusOK,
			`{"status":"ok"}`,
		)
	}
}

func writeHealthResponse(
	w http.ResponseWriter,
	statusCode int,
	body string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(statusCode)

	_, _ = w.Write([]byte(body))
}
