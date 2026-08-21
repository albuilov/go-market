package http

import (
	"net/http"

	"go-market/gen/openapi"

	swaggerui "github.com/swaggest/swgui/v5"
)

const (
	swaggerBasePath = "/docs/"
	openAPIPath     = "/openapi.json"
)

// NewSwaggerHandler создает обработчик документации и отладочного HTTP API.
// Отладочный сервер также обслуживает API-маршруты, чтобы Swagger UI мог
// выполнять запросы к ним без отдельной настройки CORS.
func NewSwaggerHandler(apiHandler http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.Handle(swaggerBasePath, swaggerui.New(
		"Go Market API",
		openAPIPath,
		swaggerBasePath,
	))
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, swaggerBasePath, http.StatusPermanentRedirect)
	})
	mux.HandleFunc("GET "+openAPIPath, openAPIHandler)
	mux.Handle("/", apiHandler)

	return mux
}

func openAPIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")

	_, _ = w.Write(openapi.Specification)
}
