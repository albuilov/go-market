package grpc

import "context"

import catalogv1 "go-market/gen/go/catalog/v1"

const minorUnitsPerUnit = 100

type Handler struct {
	catalogv1.UnimplementedCatalogServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ListProducts(
	context.Context,
	*catalogv1.ListProductsRequest,
) (*catalogv1.ListProductsResponse, error) {
	return &catalogv1.ListProductsResponse{
		Products: []*catalogv1.Product{
			newProduct("product-1", "Mechanical Keyboard", 129900, "RUB"),
			newProduct("product-2", "Wireless Mouse", 79900, "RUB"),
			newProduct("product-3", "27-inch Monitor", 2499900, "RUB"),
		},
	}, nil
}

func newProduct(
	id string,
	name string,
	priceMinorUnits int64,
	currencyCode string,
) *catalogv1.Product {
	return &catalogv1.Product{
		Id:           id,
		Name:         name,
		Price:        float64(priceMinorUnits) / minorUnitsPerUnit,
		CurrencyCode: currencyCode,
	}
}
