package grpc

import (
	"context"
	"testing"

	catalogv1 "go-market/gen/go/catalog/v1"
)

func TestHandlerListProducts(t *testing.T) {
	handler := NewHandler()

	response, err := handler.ListProducts(
		context.Background(),
		&catalogv1.ListProductsRequest{},
	)
	if err != nil {
		t.Fatalf("ListProducts returned an error: %v", err)
	}

	if got, want := len(response.GetProducts()), 3; got != want {
		t.Fatalf("product count = %d, want %d", got, want)
	}

	if got, want := response.GetProducts()[0].GetId(), "product-1"; got != want {
		t.Errorf("first product ID = %q, want %q", got, want)
	}
}
