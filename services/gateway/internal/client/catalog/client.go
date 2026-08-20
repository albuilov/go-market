package catalog

import (
	catalogv1 "go-market/gen/go/catalog/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	catalogv1.CatalogServiceClient

	connection *grpc.ClientConn
}

func New(address string) (*Client, error) {
	connection, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		CatalogServiceClient: catalogv1.NewCatalogServiceClient(connection),
		connection:           connection,
	}, nil
}

func (c *Client) Close() error {
	return c.connection.Close()
}
