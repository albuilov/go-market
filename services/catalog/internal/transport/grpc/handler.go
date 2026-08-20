package grpc

import "context"

import catalogv1 "go-market/gen/go/catalog/v1"

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
			{
				Id:              "product-1",
				Name:            "Mechanical Keyboard",
				PriceMinorUnits: 129900,
				CurrencyCode:    "RUB",
			},
			{
				Id:              "product-2",
				Name:            "Wireless Mouse",
				PriceMinorUnits: 79900,
				CurrencyCode:    "RUB",
			},
			{
				Id:              "product-3",
				Name:            "27-inch Monitor",
				PriceMinorUnits: 2499900,
				CurrencyCode:    "RUB",
			},
		},
	}, nil
}
