package app

import (
	"context"
	"errors"
	"fmt"
	"go-market/pkg/grpcmiddleware"
	"log/slog"
	"net"
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
		grpcmiddleware.TimeoutUnaryClientInterceptor(cfg.Catalog.GRPCRequestTimeout),
		authInterceptor.UnaryClient,
	)
	if err != nil {
		return fmt.Errorf("create catalog client: %w", err)
	}
	defer catalog.Close()

	handler, err := transporthttp.NewHandler(
		ctx,
		logger,
		catalog,
		catalog.HealthClient,
	)
	if err != nil {
		return fmt.Errorf("create gateway HTTP handler: %w", err)
	}

	listener, err := net.Listen(
		"tcp",
		cfg.Gateway.HTTPAddress,
	)
	if err != nil {
		return fmt.Errorf(
			"listen for gateway HTTP connections: %w",
			err,
		)
	}
	defer listener.Close()

	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: cfg.Gateway.ReadHeaderTimeout,
		ReadTimeout:       cfg.Gateway.ReadTimeout,
		WriteTimeout:      cfg.Gateway.WriteTimeout,
		IdleTimeout:       cfg.Gateway.IdleTimeout,
	}

	serveError := make(chan error, 1)

	go func() {
		serveError <- server.Serve(listener)
	}()

	select {
	case err := <-serveError:
		if err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf(
				"serve gateway HTTP requests: %w",
				err,
			)
		}

		return nil

	case <-ctx.Done():
	}

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		cfg.Gateway.ShutdownTimeout,
	)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		logger.Warn(
			"gateway graceful shutdown failed, forcing close",
			slog.Duration(
				"timeout",
				cfg.Gateway.ShutdownTimeout,
			),
			slog.Any("error", err),
		)

		if closeErr := server.Close(); closeErr != nil &&
			!errors.Is(closeErr, http.ErrServerClosed) {
			return fmt.Errorf(
				"force close gateway HTTP server: %w",
				closeErr,
			)
		}
	}

	if err := <-serveError; err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf(
			"serve gateway HTTP requests: %w",
			err,
		)
	}

	return nil
}
