package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"go-market/catalog/internal/config"
	transportgrpc "go-market/catalog/internal/transport/grpc"
	catalogv1 "go-market/gen/go/catalog/v1"
	"go-market/pkg/grpcmiddleware"

	"buf.build/go/protovalidate"
	grpcprotovalidate "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
)

func Run(
	ctx context.Context,
	logger *slog.Logger,
	cfg config.Config,
) error {
	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("create protovalidate validator: %w", err)
	}

	listener, err := net.Listen("tcp", cfg.Catalog.GRPCAddress)
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.RequestIDUnaryServerInterceptor,
			grpcmiddleware.LoggingUnaryServerInterceptor(logger),
			grpcmiddleware.RecoveryUnaryServerInterceptor(logger),
			grpcprotovalidate.UnaryServerInterceptor(validator),
		),
	)
	handler := transportgrpc.NewHandler()

	catalogv1.RegisterCatalogServiceServer(server, handler)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	return server.Serve(listener)
}
