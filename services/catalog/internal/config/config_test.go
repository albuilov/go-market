package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("CATALOG_GRPC_ADDRESS", "  127.0.0.1:50052  ")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Catalog.GRPCAddress != "127.0.0.1:50052" {
		t.Errorf(
			"Catalog.GRPCAddress = %q, want %q",
			cfg.Catalog.GRPCAddress,
			"127.0.0.1:50052",
		)
	}
}

func TestLoadUsesDefaultCatalogGRPCAddress(t *testing.T) {
	t.Setenv("CATALOG_GRPC_ADDRESS", "  ")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Catalog.GRPCAddress != defaultCatalogGRPCAddress {
		t.Errorf(
			"Catalog.GRPCAddress = %q, want %q",
			cfg.Catalog.GRPCAddress,
			defaultCatalogGRPCAddress,
		)
	}
}
