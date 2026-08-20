package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	catalogv1 "go-market/gen/go/catalog/v1"

	"google.golang.org/grpc"
)

type catalogClientStub struct {
	response *catalogv1.ListProductsResponse
	err      error
}

func (s catalogClientStub) ListProducts(
	context.Context,
	*catalogv1.ListProductsRequest,
	...grpc.CallOption,
) (*catalogv1.ListProductsResponse, error) {
	return s.response, s.err
}

func TestHandlerListProducts(t *testing.T) {
	handler := NewHandler(catalogClientStub{
		response: &catalogv1.ListProductsResponse{
			Products: []*catalogv1.Product{
				{
					Id:              "product-1",
					Name:            "Keyboard",
					PriceMinorUnits: 129900,
					CurrencyCode:    "RUB",
				},
			},
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/products", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if got, want := recorder.Code, http.StatusOK; got != want {
		t.Fatalf("status code = %d, want %d", got, want)
	}

	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	if !strings.Contains(recorder.Body.String(), `"id":"product-1"`) {
		t.Errorf("response body does not contain product: %s", recorder.Body.String())
	}
}
