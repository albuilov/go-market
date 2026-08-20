package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	catalogv1 "go-market/gen/go/catalog/v1"

	"google.golang.org/grpc"
)

type catalogClient interface {
	ListProducts(
		ctx context.Context,
		in *catalogv1.ListProductsRequest,
		opts ...grpc.CallOption,
	) (*catalogv1.ListProductsResponse, error)
}

type Handler struct {
	catalog catalogClient
}

type listProductsResponse struct {
	Products []productResponse `json:"products"`
}

type productResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	PriceMinorUnits int64  `json:"price_minor_units"`
	CurrencyCode    string `json:"currency_code"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(catalog catalogClient) http.Handler {
	handler := &Handler{catalog: catalog}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", handler.listProducts)

	return mux
}

func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	response, err := h.catalog.ListProducts(
		r.Context(),
		&catalogv1.ListProductsRequest{},
	)
	if err != nil {
		log.Printf("list products through catalog service: %v", err)
		writeJSON(w, http.StatusBadGateway, errorResponse{
			Error: "catalog service is unavailable",
		})
		return
	}

	products := make([]productResponse, 0, len(response.GetProducts()))
	for _, product := range response.GetProducts() {
		products = append(products, productResponse{
			ID:              product.GetId(),
			Name:            product.GetName(),
			PriceMinorUnits: product.GetPriceMinorUnits(),
			CurrencyCode:    product.GetCurrencyCode(),
		})
	}

	writeJSON(w, http.StatusOK, listProductsResponse{Products: products})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode HTTP response: %v", err)
	}
}
