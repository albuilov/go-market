package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"go-market/pkg/grpcmiddleware"

	gatewayauth "go-market/gateway/internal/auth"
	catalogclient "go-market/gateway/internal/client/catalog"
	"go-market/gateway/internal/config"
	transporthttp "go-market/gateway/internal/transport/http"
)

type namedHTTPServer struct {
	name     string
	server   *http.Server
	listener net.Listener
}

type httpServeResult struct {
	name string
	err  error
}

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

	swaggerHandler := transporthttp.NewSwaggerHandler(handler)

	apiListener, err := net.Listen("tcp", cfg.Gateway.HTTPAddress)
	if err != nil {
		return fmt.Errorf("listen for gateway HTTP connections: %w", err)
	}
	defer apiListener.Close()

	swaggerListener, err := net.Listen("tcp", cfg.Gateway.SwaggerHTTPAddress)
	if err != nil {
		return fmt.Errorf("listen for gateway Swagger HTTP connections: %w", err)
	}
	defer swaggerListener.Close()

	servers := []namedHTTPServer{
		{
			name: "api",
			server: &http.Server{
				Handler:           handler,
				ReadHeaderTimeout: cfg.Gateway.ReadHeaderTimeout,
				ReadTimeout:       cfg.Gateway.ReadTimeout,
				WriteTimeout:      cfg.Gateway.WriteTimeout,
				IdleTimeout:       cfg.Gateway.IdleTimeout,
			},
			listener: apiListener,
		},
		{
			name: "swagger",
			server: &http.Server{
				Handler:           swaggerHandler,
				ReadHeaderTimeout: cfg.Gateway.ReadHeaderTimeout,
				ReadTimeout:       cfg.Gateway.ReadTimeout,
				WriteTimeout:      cfg.Gateway.WriteTimeout,
				IdleTimeout:       cfg.Gateway.IdleTimeout,
			},
			listener: swaggerListener,
		},
	}

	serveResults := make(chan httpServeResult, len(servers))

	for _, namedServer := range servers {
		go func() {
			serveResults <- httpServeResult{
				name: namedServer.name,
				err:  namedServer.server.Serve(namedServer.listener),
			}
		}()
	}

	remainingServers := len(servers)
	var runError error

	select {
	case result := <-serveResults:
		remainingServers--
		runError = unexpectedServeError(result)
	case <-ctx.Done():
	}

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		cfg.Gateway.ShutdownTimeout,
	)
	defer cancel()

	var shutdownErrors []error

	for _, namedServer := range servers {
		if err := namedServer.server.Shutdown(shutdownContext); err != nil {
			logger.Warn(
				"gateway graceful shutdown failed, forcing close",
				slog.String("server", namedServer.name),
				slog.Duration("timeout", cfg.Gateway.ShutdownTimeout),
				slog.Any("error", err),
			)

			if closeErr := namedServer.server.Close(); closeErr != nil &&
				!errors.Is(closeErr, http.ErrServerClosed) {
				shutdownErrors = append(
					shutdownErrors,
					fmt.Errorf("force close gateway %s HTTP server: %w", namedServer.name, closeErr),
				)
			}
		}
	}

	for range remainingServers {
		result := <-serveResults
		if err := unexpectedServeError(result); err != nil {
			shutdownErrors = append(shutdownErrors, err)
		}
	}

	return errors.Join(append([]error{runError}, shutdownErrors...)...)
}

func unexpectedServeError(result httpServeResult) error {
	if result.err == nil || errors.Is(result.err, http.ErrServerClosed) {
		return nil
	}

	return fmt.Errorf("serve gateway %s HTTP requests: %w", result.name, result.err)
}
