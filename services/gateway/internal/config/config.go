package config

import (
	"fmt"
	"time"

	"go-market/pkg/envconfig"
)

const (
	defaultGatewayHTTPAddress       = ":3000"
	defaultGatewayReadHeaderTimeout = 5 * time.Second
	defaultGatewayShutdownTimeout   = 5 * time.Second
)

type Config struct {
	Gateway GatewayConfig
	Catalog CatalogConfig
	JWT     JWTConfig
}

type GatewayConfig struct {
	HTTPAddress       string
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

type CatalogConfig struct {
	GRPCAddress string
}

type JWTConfig struct {
	Secret   string
	Issuer   string
	Audience string
}

func Load() (Config, error) {
	catalogGRPCAddress, err := envconfig.Required("CATALOG_GRPC_ADDRESS")
	if err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}

	jwtSecret, err := envconfig.Required("JWT_SECRET")
	if err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}

	jwtIssuer, err := envconfig.Required("JWT_ISSUER")
	if err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}

	jwtAudience, err := envconfig.Required("JWT_AUDIENCE")
	if err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}

	return Config{
		Gateway: GatewayConfig{
			HTTPAddress:       envconfig.OrDefault("GATEWAY_HTTP_ADDRESS", defaultGatewayHTTPAddress),
			ReadHeaderTimeout: defaultGatewayReadHeaderTimeout,
			ShutdownTimeout:   defaultGatewayShutdownTimeout,
		},
		Catalog: CatalogConfig{
			GRPCAddress: catalogGRPCAddress,
		},
		JWT: JWTConfig{
			Secret:   jwtSecret,
			Issuer:   jwtIssuer,
			Audience: jwtAudience,
		},
	}, nil
}
