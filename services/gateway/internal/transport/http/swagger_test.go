package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAPIProductMatchesPublicJSON(t *testing.T) {
	var specification struct {
		Definitions map[string]struct {
			Properties map[string]json.RawMessage `json:"properties"`
		} `json:"definitions"`
	}

	if err := json.Unmarshal(openapiSpecification(), &specification); err != nil {
		t.Fatalf("unmarshal OpenAPI specification: %v", err)
	}

	properties := specification.Definitions["catalog.v1.Product"].Properties
	for _, field := range []string{"id", "name", "price", "currency_code"} {
		if _, ok := properties[field]; !ok {
			t.Errorf("OpenAPI Product does not contain %q", field)
		}
	}
	for _, internalField := range []string{"price_minor_units", "priceMinorUnits"} {
		if _, ok := properties[internalField]; ok {
			t.Errorf("OpenAPI Product exposes internal field %q", internalField)
		}
	}
}

func openapiSpecification() []byte {
	handler := NewSwaggerHandler(http.NotFoundHandler())
	request := httptest.NewRequest(http.MethodGet, openAPIPath, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	return response.Body.Bytes()
}

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
