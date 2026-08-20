package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go-market/gateway/internal/app"
	gatewayauth "go-market/gateway/internal/auth"
)

func main() {
	httpAddress := os.Getenv("HTTP_ADDRESS")
	if httpAddress == "" {
		httpAddress = ":8080"
	}

	catalogAddress := os.Getenv("CATALOG_GRPC_ADDRESS")
	if catalogAddress == "" {
		catalogAddress = "localhost:50051"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtIssuer := os.Getenv("JWT_ISSUER")
	jwtAudience := os.Getenv("JWT_AUDIENCE")

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	tokenVerifier, err := gatewayauth.NewVerifier(
		jwtSecret,
		jwtIssuer,
		jwtAudience,
	)
	if err != nil {
		logger.Error(
			"failed to configure JWT verifier",
			slog.Any("error", err),
		)
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
		slog.String("http_address", httpAddress),
		slog.String("catalog_address", catalogAddress),
	)

	if err := app.Run(
		ctx,
		logger,
		httpAddress,
		catalogAddress,
		tokenVerifier,
	); err != nil {
		logger.Error(
			"gateway service failed",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}
