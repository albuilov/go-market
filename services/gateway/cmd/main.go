package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-market/gateway/internal/app"
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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	log.Printf(
		"gateway HTTP server is listening on %s; catalog address is %s",
		httpAddress,
		catalogAddress,
	)

	if err := app.Run(ctx, httpAddress, catalogAddress); err != nil {
		log.Fatalf("gateway service failed: %v", err)
	}
}
