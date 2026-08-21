package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-market/catalog/internal/app"
	"go-market/catalog/internal/config"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	log.Printf(
		"catalog gRPC server is listening on %s",
		cfg.Catalog.GRPCAddress,
	)

	if err := app.Run(ctx, cfg); err != nil &&
		!errors.Is(err, grpc.ErrServerStopped) {
		log.Fatalf("catalog service failed: %v", err)
	}
}
