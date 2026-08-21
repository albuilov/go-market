package app

import (
	"context"
	"net"

	"go-market/catalog/internal/config"
	transportgrpc "go-market/catalog/internal/transport/grpc"
	catalogv1 "go-market/gen/go/catalog/v1"

	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg config.Config) error {
	listener, err := net.Listen("tcp", cfg.Catalog.GRPCAddress)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	handler := transportgrpc.NewHandler()

	catalogv1.RegisterCatalogServiceServer(server, handler)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	return server.Serve(listener)
}
