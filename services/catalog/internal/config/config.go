package config

import (
	"time"

	"go-market/pkg/envconfig"
)

const (
	defaultCatalogGRPCAddress     = ":50051"
	defaultCatalogShutdownTimeout = 10 * time.Second
)

type Config struct {
	Catalog CatalogConfig
}

type CatalogConfig struct {
	GRPCAddress     string
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	return Config{
		Catalog: CatalogConfig{
			GRPCAddress:     envconfig.OrDefault("CATALOG_GRPC_ADDRESS", defaultCatalogGRPCAddress),
			ShutdownTimeout: defaultCatalogShutdownTimeout,
		},
	}, nil
}
