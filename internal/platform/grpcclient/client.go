// Package grpcclient создает стандартные соединения с gRPC-сервисами.
package grpcclient

import (
	"errors"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Config содержит обязательные параметры gRPC-соединения.
type Config struct {
	Address              string
	RequestTimeout       time.Duration
	TransportCredentials credentials.TransportCredentials
}

// New создает соединение с общим timeout interceptor.
// Дополнительные interceptors выполняются внутри общего timeout.
func New(
	config Config,
	additionalInterceptors ...grpc.UnaryClientInterceptor,
) (*grpc.ClientConn, error) {
	address := strings.TrimSpace(config.Address)
	if address == "" {
		return nil, errors.New("gRPC address is required")
	}
	if config.RequestTimeout <= 0 {
		return nil, errors.New("gRPC request timeout must be positive")
	}
	if config.TransportCredentials == nil {
		return nil, errors.New("gRPC transport credentials are required")
	}

	interceptors := []grpc.UnaryClientInterceptor{
		TimeoutUnaryClientInterceptor(config.RequestTimeout),
	}
	interceptors = append(interceptors, additionalInterceptors...)

	return grpc.NewClient(
		address,
		grpc.WithTransportCredentials(config.TransportCredentials),
		grpc.WithChainUnaryInterceptor(interceptors...),
	)
}
