package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type requestStateKey struct{}

type requestState struct {
	err error
}

// RecordError сохраняет ошибку для единственной итоговой записи HTTP-запроса.
func RecordError(ctx context.Context, err error) {
	state, ok := ctx.Value(requestStateKey{}).(*requestState)
	if !ok {
		return
	}

	state.err = err
}

type responseWriter struct {
	http.ResponseWriter

	status int
	bytes  int
}

func (w *responseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}

	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}

	n, err := w.ResponseWriter.Write(body)
	w.bytes += n

	return n, err
}

func (w *responseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *responseWriter) Flush() {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}

	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func Logging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		writer := &responseWriter{ResponseWriter: w}
		state := &requestState{}
		r = r.WithContext(
			context.WithValue(r.Context(), requestStateKey{}, state),
		)

		next.ServeHTTP(writer, r)

		status := writer.status
		if status == 0 {
			status = http.StatusOK
		}
		if isHealthEndpoint(r.URL.Path) && status < http.StatusBadRequest {
			return
		}

		level := slog.LevelInfo
		switch {
		case status >= http.StatusInternalServerError:
			level = slog.LevelError
		case status >= http.StatusBadRequest:
			level = slog.LevelWarn
		}

		requestID, _ := RequestIDFromContext(r.Context())

		attributes := []any{
			slog.String("request_id", requestID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", status),
			slog.Int("response_bytes", writer.bytes),
			slog.Duration("duration", time.Since(startedAt)),
		}
		if state.err != nil {
			attributes = append(attributes, slog.Any("error", state.err))
		}

		logger.Log(
			r.Context(),
			level,
			"HTTP request completed",
			attributes...,
		)
	})
}

func isHealthEndpoint(path string) bool {
	return path == "/healthz" || path == "/readyz"
}
