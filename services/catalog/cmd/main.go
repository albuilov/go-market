package main

import (
	"context"
	"errors"
	"go-market/catalog/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	address := os.Getenv("GRPC_ADDRESS")
	if address == "" {
		address = ":50051"
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	log.Printf("catalog gRPC server is listening on %s", address)

	if err := app.Run(ctx, address); err != nil &&
		!errors.Is(err, grpc.ErrServerStopped) {
		log.Fatalf("catalog service failed: %v", err)
	}
}
