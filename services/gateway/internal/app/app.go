package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	catalogclient "go-market/gateway/internal/client/catalog"
	transporthttp "go-market/gateway/internal/transport/http"
)

const shutdownTimeout = 5 * time.Second

func Run(ctx context.Context, httpAddress, catalogAddress string) error {
	catalog, err := catalogclient.New(catalogAddress)
	if err != nil {
		return err
	}
	defer catalog.Close()

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           transporthttp.NewHandler(catalog),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownContext, cancel := context.WithTimeout(
			context.Background(),
			shutdownTimeout,
		)
		defer cancel()

		_ = server.Shutdown(shutdownContext)
	}()

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
