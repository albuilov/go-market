package config

import "go-market/pkg/envconfig"

const defaultCatalogGRPCAddress = ":50051"

type Config struct {
	Catalog CatalogConfig
}

type CatalogConfig struct {
	GRPCAddress string
}

func Load() (Config, error) {
	return Config{
		Catalog: CatalogConfig{
			GRPCAddress: envconfig.OrDefault(
				"CATALOG_GRPC_ADDRESS",
				defaultCatalogGRPCAddress,
			),
		},
	}, nil
}
