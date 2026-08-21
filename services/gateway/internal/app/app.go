package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	gatewayauth "go-market/gateway/internal/auth"
	catalogclient "go-market/gateway/internal/client/catalog"
	"go-market/gateway/internal/config"
	transporthttp "go-market/gateway/internal/transport/http"
)

func Run(
	ctx context.Context,
	logger *slog.Logger,
	cfg config.Config,
	tokenVerifier gatewayauth.TokenVerifier,
) error {
	authInterceptor := gatewayauth.NewInterceptor(
		logger,
		tokenVerifier,
	)

	catalog, err := catalogclient.New(
		cfg.Catalog.GRPCAddress,
		authInterceptor.UnaryClient,
	)
	if err != nil {
		return err
	}
	defer catalog.Close()

	handler, err := transporthttp.NewHandler(ctx, logger, catalog, catalog.HealthClient)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              cfg.Gateway.HTTPAddress,
		Handler:           handler,
		ReadHeaderTimeout: cfg.Gateway.ReadHeaderTimeout,
	}

	go func() {
		<-ctx.Done()

		shutdownContext, cancel := context.WithTimeout(
			context.Background(),
			cfg.Gateway.ShutdownTimeout,
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
