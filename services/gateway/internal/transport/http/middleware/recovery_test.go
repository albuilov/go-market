package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoveryHandlesPanic(t *testing.T) {
	var logs bytes.Buffer

	logger := slog.New(
		slog.NewJSONHandler(&logs, nil),
	)

	handler := RequestID(
		Logging(
			logger,
			Recovery(
				logger,
				http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
					panic("unexpected failure")
				}),
			),
		),
	)

	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if got, want := recorder.Code, http.StatusInternalServerError; got != want {
		t.Fatalf("status code = %d, want %d", got, want)
	}

	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	if got := recorder.Header().Get(RequestIDHeader); got == "" {
		t.Error("response does not contain request ID")
	}

	if !strings.Contains(
		recorder.Body.String(),
		`"message":"internal server error"`,
	) {
		t.Errorf("unexpected response body: %s", recorder.Body.String())
	}

	if !strings.Contains(logs.String(), `"msg":"panic recovered"`) {
		t.Errorf("panic was not logged: %s", logs.String())
	}
}
