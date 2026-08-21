package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go-market/gateway/internal/app"
	gatewayauth "go-market/gateway/internal/auth"
	"go-market/gateway/internal/config"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	tokenVerifier, err := gatewayauth.NewVerifier(
		cfg.JWT.Secret,
		cfg.JWT.Issuer,
		cfg.JWT.Audience,
	)
	if err != nil {
		logger.Error("failed to configure JWT verifier", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	logger.Info(
		"starting gateway service",
		slog.String("http_address", cfg.Gateway.HTTPAddress),
		slog.String("swagger_http_address", cfg.Gateway.SwaggerHTTPAddress),
		slog.String("catalog_address", cfg.Catalog.GRPCAddress),
	)

	if err := app.Run(
		ctx,
		logger,
		cfg,
		tokenVerifier,
	); err != nil {
		logger.Error("gateway service failed", slog.Any("error", err))
		os.Exit(1)
	}
}
