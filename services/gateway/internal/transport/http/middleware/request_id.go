package middleware

import (
	"context"
	"net/http"

	"go-market/pkg/requestid"
)

// RequestIDHeader содержит имя HTTP-заголовка с идентификатором запроса.
const RequestIDHeader = requestid.HTTPHeader

// RequestID добавляет идентификатор запроса в HTTP-заголовки и context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := requestid.Normalize(
			r.Header.Get(RequestIDHeader),
		)

		r.Header.Set(RequestIDHeader, id)
		w.Header().Set(RequestIDHeader, id)

		ctx := requestid.WithID(r.Context(), id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext возвращает идентификатор запроса из context.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	return requestid.FromContext(ctx)
}
