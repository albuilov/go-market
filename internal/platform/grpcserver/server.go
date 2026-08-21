// Package grpcserver собирает стандартный gRPC-сервер приложения.
package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"buf.build/go/protovalidate"
	grpcprotovalidate "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

// Server объединяет gRPC-сервер и стандартный Health API.
type Server struct {
	logger       *slog.Logger
	grpcServer   *grpc.Server
	healthServer *health.Server
}

// New создает сервер со стандартными unary interceptors.
// Дополнительные interceptors выполняются перед Protovalidate.
func New(
	logger *slog.Logger,
	additionalInterceptors ...grpc.UnaryServerInterceptor,
) (*Server, error) {
	// TODO: добавить такую же стандартную цепочку для stream interceptors
	// TODO: перед появлением первого streaming RPC.
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("create protovalidate validator: %w", err)
	}

	interceptors := []grpc.UnaryServerInterceptor{
		RequestIDUnaryServerInterceptor,
		LoggingUnaryServerInterceptor(logger),
		RecoveryUnaryServerInterceptor(logger),
	}
	interceptors = append(interceptors, additionalInterceptors...)
	interceptors = append(
		interceptors,
		grpcprotovalidate.UnaryServerInterceptor(validator),
	)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)
	healthServer := health.NewServer()
	healthv1.RegisterHealthServer(grpcServer, healthServer)

	return &Server{
		logger:       logger,
		grpcServer:   grpcServer,
		healthServer: healthServer,
	}, nil
}

// GRPC возвращает сервер для регистрации бизнес-сервисов.
func (s *Server) GRPC() *grpc.Server {
	return s.grpcServer
}

// SetServing отмечает зарегистрированные сервисы готовыми принимать запросы.
func (s *Server) SetServing(services ...string) {
	for _, service := range services {
		s.healthServer.SetServingStatus(
			service,
			healthv1.HealthCheckResponse_SERVING,
		)
	}
}

// Run запускает сервер и ограничивает время graceful shutdown после отмены context.
func (s *Server) Run(
	ctx context.Context,
	address string,
	shutdownTimeout time.Duration,
) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen for gRPC connections: %w", err)
	}
	defer listener.Close()

	serveError := make(chan error, 1)

	go func() {
		serveError <- s.grpcServer.Serve(listener)
	}()

	select {
	case err := <-serveError:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("serve gRPC requests: %w", err)
		}

		return nil

	case <-ctx.Done():
	}

	s.healthServer.Shutdown()

	gracefulStopDone := make(chan struct{})

	go func() {
		s.grpcServer.GracefulStop()
		close(gracefulStopDone)
	}()

	timer := time.NewTimer(shutdownTimeout)
	defer timer.Stop()

	select {
	case <-gracefulStopDone:

	case <-timer.C:
		s.logger.Warn(
			"gRPC graceful shutdown timed out",
			slog.Duration("timeout", shutdownTimeout),
		)

		s.grpcServer.Stop()
		<-gracefulStopDone
	}

	if err := <-serveError; err != nil &&
		!errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("serve gRPC requests: %w", err)
	}

	return nil
}
