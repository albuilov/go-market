package app

import (
	"context"
	"fmt"
	"log/slog"

	catalogv1 "go-market/gen/go/catalog/v1"
	sharedgrpc "go-market/internal/platform/grpcserver"
	"go-market/services/catalog/internal/config"
	transportgrpc "go-market/services/catalog/internal/transport/grpc"
)

func Run(
	ctx context.Context,
	logger *slog.Logger,
	cfg config.Config,
) error {
	server, err := sharedgrpc.New(logger)
	if err != nil {
		return fmt.Errorf("create catalog gRPC server: %w", err)
	}

	handler := transportgrpc.NewHandler()
	catalogv1.RegisterCatalogServiceServer(server.GRPC(), handler)

	server.SetServing(
		catalogv1.CatalogService_ServiceDesc.ServiceName,
	)

	if err := server.Run(
		ctx,
		cfg.Catalog.GRPCAddress,
		cfg.Catalog.ShutdownTimeout,
	); err != nil {
		return fmt.Errorf("run catalog gRPC server: %w", err)
	}

	return nil
}
