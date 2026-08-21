package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Setenv("GATEWAY_HTTP_ADDRESS", "  :9090  ")
	t.Setenv("SWAGGER_HTTP_ADDRESS", "  :9091  ")
	t.Setenv("CATALOG_GRPC_ADDRESS", "  catalog:50051  ")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("JWT_ISSUER", "issuer")
	t.Setenv("JWT_AUDIENCE", "audience")

	cfg, err := Load()

	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Gateway.HTTPAddress != ":9090" {
		t.Errorf("Gateway.HTTPAddress = %q, want %q", cfg.Gateway.HTTPAddress, ":9090")
	}

	if cfg.Gateway.SwaggerHTTPAddress != ":9091" {
		t.Errorf(
			"Gateway.SwaggerHTTPAddress = %q, want %q",
			cfg.Gateway.SwaggerHTTPAddress,
			":9091",
		)
	}

	if cfg.Gateway.ReadHeaderTimeout != 5*time.Second {
		t.Errorf(
			"Gateway.ReadHeaderTimeout = %s, want %s",
			cfg.Gateway.ReadHeaderTimeout,
			5*time.Second,
		)
	}

	if cfg.Gateway.ShutdownTimeout != 5*time.Second {
		t.Errorf("Gateway.ShutdownTimeout = %s, want %s", cfg.Gateway.ShutdownTimeout, 5*time.Second)
	}

	if cfg.Catalog.GRPCAddress != "catalog:50051" {
		t.Errorf("Catalog.GRPCAddress = %q, want %q", cfg.Catalog.GRPCAddress, "catalog:50051")
	}

	if cfg.Gateway.ReadTimeout != defaultGatewayReadTimeout {
		t.Errorf(
			"Gateway.ReadTimeout = %s, want %s",
			cfg.Gateway.ReadTimeout,
			defaultGatewayReadTimeout,
		)
	}

	if cfg.Gateway.WriteTimeout != defaultGatewayWriteTimeout {
		t.Errorf(
			"Gateway.WriteTimeout = %s, want %s",
			cfg.Gateway.WriteTimeout,
			defaultGatewayWriteTimeout,
		)
	}

	if cfg.Gateway.IdleTimeout != defaultGatewayIdleTimeout {
		t.Errorf(
			"Gateway.IdleTimeout = %s, want %s",
			cfg.Gateway.IdleTimeout,
			defaultGatewayIdleTimeout,
		)
	}

	if cfg.Catalog.GRPCRequestTimeout != defaultCatalogGRPCRequestTimeout {
		t.Errorf(
			"Catalog.GRPCRequestTimeout = %s, want %s",
			cfg.Catalog.GRPCRequestTimeout,
			defaultCatalogGRPCRequestTimeout,
		)
	}
}

func TestLoadUsesDefaultGatewayHTTPAddress(t *testing.T) {
	t.Setenv("GATEWAY_HTTP_ADDRESS", "  ")
	t.Setenv("SWAGGER_HTTP_ADDRESS", "  ")
	setRequiredEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Gateway.HTTPAddress != defaultGatewayHTTPAddress {
		t.Errorf(
			"Gateway.HTTPAddress = %q, want %q",
			cfg.Gateway.HTTPAddress,
			defaultGatewayHTTPAddress,
		)
	}

	if cfg.Gateway.SwaggerHTTPAddress != defaultGatewaySwaggerHTTPAddress {
		t.Errorf(
			"Gateway.SwaggerHTTPAddress = %q, want %q",
			cfg.Gateway.SwaggerHTTPAddress,
			defaultGatewaySwaggerHTTPAddress,
		)
	}
}

func TestLoadRejectsEmptyRequiredEnvironmentVariables(t *testing.T) {
	requiredKeys := []string{
		"CATALOG_GRPC_ADDRESS",
		"JWT_SECRET",
		"JWT_ISSUER",
		"JWT_AUDIENCE",
	}

	for _, key := range requiredKeys {
		t.Run(key, func(t *testing.T) {
			setRequiredEnv(t)
			t.Setenv(key, "  ")

			_, err := Load()
			if err == nil {
				t.Fatal("Load() error = nil, want an error")
			}
			if !strings.Contains(err.Error(), key) {
				t.Errorf("Load() error = %q, want it to contain %q", err, key)
			}
		})
	}
}

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("CATALOG_GRPC_ADDRESS", "catalog:50051")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("JWT_ISSUER", "issuer")
	t.Setenv("JWT_AUDIENCE", "audience")
}
