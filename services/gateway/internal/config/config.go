package config

import (
	"fmt"
	"time"

	"go-market/pkg/envconfig"
)

const (
	defaultGatewayHTTPAddress        = ":3000"
	defaultGatewaySwaggerHTTPAddress = ":3001"
	defaultGatewayReadHeaderTimeout  = 5 * time.Second
	defaultGatewayReadTimeout        = 15 * time.Second
	defaultGatewayWriteTimeout       = 15 * time.Second
	defaultGatewayIdleTimeout        = 60 * time.Second
	defaultGatewayShutdownTimeout    = 5 * time.Second
	defaultCatalogGRPCRequestTimeout = 10 * time.Second
)

type Config struct {
	Gateway GatewayConfig
	Catalog CatalogConfig
	JWT     JWTConfig
}

type GatewayConfig struct {
	HTTPAddress        string
	SwaggerHTTPAddress string
	ReadHeaderTimeout  time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	ShutdownTimeout    time.Duration
}

type CatalogConfig struct {
	GRPCAddress        string
	GRPCRequestTimeout time.Duration
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
			HTTPAddress:        envconfig.OrDefault("GATEWAY_HTTP_ADDRESS", defaultGatewayHTTPAddress),
			SwaggerHTTPAddress: envconfig.OrDefault("SWAGGER_HTTP_ADDRESS", defaultGatewaySwaggerHTTPAddress),
			ReadHeaderTimeout:  defaultGatewayReadHeaderTimeout,
			ReadTimeout:        defaultGatewayReadTimeout,
			WriteTimeout:       defaultGatewayWriteTimeout,
			IdleTimeout:        defaultGatewayIdleTimeout,
			ShutdownTimeout:    defaultGatewayShutdownTimeout,
		},
		Catalog: CatalogConfig{
			GRPCAddress:        catalogGRPCAddress,
			GRPCRequestTimeout: defaultCatalogGRPCRequestTimeout,
		},
		JWT: JWTConfig{
			Secret:   jwtSecret,
			Issuer:   jwtIssuer,
			Audience: jwtAudience,
		},
	}, nil
}
