package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	gatewayauth "go-market/gateway/internal/auth"
	catalogclient "go-market/gateway/internal/client/catalog"
	transporthttp "go-market/gateway/internal/transport/http"
)

const shutdownTimeout = 5 * time.Second

func Run(
	ctx context.Context,
	logger *slog.Logger,
	httpAddress string,
	catalogAddress string,
	tokenVerifier gatewayauth.TokenVerifier,
) error {
	authInterceptor := gatewayauth.NewInterceptor(
		logger,
		tokenVerifier,
	)

	catalog, err := catalogclient.New(catalogAddress, authInterceptor.UnaryClient)
	if err != nil {
		return err
	}
	defer catalog.Close()

	handler, err := transporthttp.NewHandler(ctx, logger, catalog)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           handler,
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
