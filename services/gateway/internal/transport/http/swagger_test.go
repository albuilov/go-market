package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSwaggerHandlerServesOpenAPISpecification(t *testing.T) {
	handler := NewSwaggerHandler(http.NotFoundHandler())
	request := httptest.NewRequest(http.MethodGet, openAPIPath, nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusOK)
	}

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json; charset=utf-8")
	}

	if !strings.Contains(response.Body.String(), `"/api/v1/products"`) {
		t.Error("OpenAPI specification does not contain the products endpoint")
	}
}

func TestSwaggerHandlerServesUI(t *testing.T) {
	handler := NewSwaggerHandler(http.NotFoundHandler())
	request := httptest.NewRequest(http.MethodGet, swaggerBasePath, nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusOK)
	}

	if !strings.Contains(response.Body.String(), "Go Market API") {
		t.Error("Swagger UI response does not contain the API title")
	}
}

func TestSwaggerHandlerRedirectsToUI(t *testing.T) {
	handler := NewSwaggerHandler(http.NotFoundHandler())
	request := httptest.NewRequest(http.MethodGet, "/docs", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusPermanentRedirect {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusPermanentRedirect)
	}

	if location := response.Header().Get("Location"); location != swaggerBasePath {
		t.Errorf("Location = %q, want %q", location, swaggerBasePath)
	}
}

func TestSwaggerHandlerForwardsAPIRequests(t *testing.T) {
	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler := NewSwaggerHandler(apiHandler)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusNoContent)
	}
}
