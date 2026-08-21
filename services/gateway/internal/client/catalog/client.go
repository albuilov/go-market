package catalog

import (
	"time"

	catalogv1 "go-market/gen/go/catalog/v1"
	"go-market/internal/platform/grpcclient"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

type Client struct {
	catalogv1.CatalogServiceClient

	HealthClient healthv1.HealthClient
	connection   *grpc.ClientConn
}

func New(
	address string,
	requestTimeout time.Duration,
	unaryInterceptors ...grpc.UnaryClientInterceptor,
) (*Client, error) {
	// В локальной Docker-сети шифрование пока отключено.
	// TODO: перед production-развертыванием передавать TLS или mTLS credentials.
	connection, err := grpcclient.New(
		grpcclient.Config{
			Address:              address,
			RequestTimeout:       requestTimeout,
			TransportCredentials: insecure.NewCredentials(),
		},
		unaryInterceptors...,
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		CatalogServiceClient: catalogv1.NewCatalogServiceClient(connection),
		HealthClient:         healthv1.NewHealthClient(connection),
		connection:           connection,
	}, nil
}

func (c *Client) Close() error {
	return c.connection.Close()
}
