package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

type requestIDContextKey struct{}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := normalizeRequestID(
			r.Header.Get(RequestIDHeader),
		)

		r.Header.Set(RequestIDHeader, requestID)
		w.Header().Set(RequestIDHeader, requestID)

		ctx := context.WithValue(
			r.Context(),
			requestIDContextKey{},
			requestID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey{}).(string)
	return requestID, ok
}

func normalizeRequestID(value string) string {
	requestID, err := uuid.Parse(value)
	if err != nil {
		return uuid.NewString()
	}

	return requestID.String()
}
