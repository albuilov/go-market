package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"go-market/catalog/internal/config"
	transportgrpc "go-market/catalog/internal/transport/grpc"
	catalogv1 "go-market/gen/go/catalog/v1"
	"go-market/pkg/grpcmiddleware"

	"buf.build/go/protovalidate"
	grpcprotovalidate "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
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

	listener, err := net.Listen(
		"tcp",
		cfg.Catalog.GRPCAddress,
	)
	if err != nil {
		return fmt.Errorf("listen for catalog gRPC connections: %w", err)
	}
	defer listener.Close()

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.RequestIDUnaryServerInterceptor,
			grpcmiddleware.LoggingUnaryServerInterceptor(logger),
			grpcmiddleware.RecoveryUnaryServerInterceptor(logger),
			grpcprotovalidate.UnaryServerInterceptor(validator),
		),
	)

	handler := transportgrpc.NewHandler()
	healthServer := health.NewServer()

	catalogv1.RegisterCatalogServiceServer(server, handler)
	healthv1.RegisterHealthServer(server, healthServer)

	healthServer.SetServingStatus(
		catalogv1.CatalogService_ServiceDesc.ServiceName,
		healthv1.HealthCheckResponse_SERVING,
	)

	serveError := make(chan error, 1)

	go func() {
		serveError <- server.Serve(listener)
	}()

	select {
	case err := <-serveError:
		if err != nil &&
			!errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf(
				"serve catalog gRPC requests: %w",
				err,
			)
		}

		return nil

	case <-ctx.Done():
	}

	healthServer.Shutdown()

	gracefulStopDone := make(chan struct{})

	go func() {
		server.GracefulStop()
		close(gracefulStopDone)
	}()

	timer := time.NewTimer(cfg.Catalog.ShutdownTimeout)
	defer timer.Stop()

	select {
	case <-gracefulStopDone:

	case <-timer.C:
		logger.Warn(
			"catalog graceful shutdown timed out",
			slog.Duration(
				"timeout",
				cfg.Catalog.ShutdownTimeout,
			),
		)

		server.Stop()
		<-gracefulStopDone
	}

	if err := <-serveError; err != nil &&
		!errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf(
			"serve catalog gRPC requests: %w",
			err,
		)
	}

	return nil
}
