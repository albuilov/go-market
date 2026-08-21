package http

import (
	"context"
	"net/http"
	"time"

	platformhealth "go-market/internal/platform/health"
	httpmiddleware "go-market/services/gateway/internal/transport/http/middleware"
)

const readinessCheckTimeout = time.Second

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
	checker platformhealth.Checker,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(
			r.Context(),
			readinessCheckTimeout,
		)
		defer cancel()

		if err := checker.Check(ctx); err != nil {
			httpmiddleware.RecordError(r.Context(), err)

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
