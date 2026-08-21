package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go-market/internal/platform/logging"
	"go-market/services/catalog/internal/app"
	"go-market/services/catalog/internal/config"

	"google.golang.org/grpc"
)

func main() {
	logger := logging.New(os.Stdout, "catalog", slog.LevelInfo)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	logger.Info(
		"starting catalog service",
		slog.String("grpc_address", cfg.Catalog.GRPCAddress),
	)

	if err := app.Run(ctx, logger, cfg); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		logger.Error("catalog service failed", slog.Any("error", err))
		os.Exit(1)
	}
}
