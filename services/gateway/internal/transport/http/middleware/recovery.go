package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type recoveryErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Recovery(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &responseWriter{ResponseWriter: w}

		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}

			requestID, _ := RequestIDFromContext(r.Context())

			logger.ErrorContext(
				r.Context(),
				"panic recovered",
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Any("panic", recovered),
				slog.String("stack", string(debug.Stack())),
			)

			// Если ответ уже начали отправлять, изменить status на 500 нельзя.
			if writer.status != 0 {
				return
			}

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusInternalServerError)

			if err := json.NewEncoder(writer).Encode(recoveryErrorResponse{
				Code:    "Internal",
				Message: "internal server error",
			}); err != nil {
				logger.ErrorContext(
					r.Context(),
					"failed to encode panic response",
					slog.Any("error", err),
				)
			}
		}()

		next.ServeHTTP(writer, r)
	})
}
