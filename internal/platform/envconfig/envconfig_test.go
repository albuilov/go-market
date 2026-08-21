package envconfig

import (
	"os"
	"strings"
	"testing"
)

func TestRequired(t *testing.T) {
	t.Setenv("TEST_REQUIRED_VALUE", "  value  ")

	value, err := Required("TEST_REQUIRED_VALUE")
	if err != nil {
		t.Fatalf("Required() error = %v", err)
	}
	if value != "value" {
		t.Errorf("Required() = %q, want %q", value, "value")
	}
}

func TestRequiredRejectsMissingAndEmptyValues(t *testing.T) {
	testCases := map[string]func(t *testing.T, key string){
		"missing": unsetEnv,
		"empty": func(t *testing.T, key string) {
			t.Setenv(key, "  ")
		},
	}

	for name, prepare := range testCases {
		t.Run(name, func(t *testing.T) {
			const key = "TEST_REQUIRED_INVALID"
			prepare(t, key)

			_, err := Required(key)
			if err == nil {
				t.Fatal("Required() error = nil, want an error")
			}
			if !strings.Contains(err.Error(), key) {
				t.Errorf("Required() error = %q, want it to contain %q", err, key)
			}
		})
	}
}

func TestOrDefault(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		t.Setenv("TEST_OPTIONAL_VALUE", "  value  ")

		if value := OrDefault("TEST_OPTIONAL_VALUE", "default"); value != "value" {
			t.Errorf("OrDefault() = %q, want %q", value, "value")
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Setenv("TEST_OPTIONAL_EMPTY", "  ")

		if value := OrDefault("TEST_OPTIONAL_EMPTY", "default"); value != "default" {
			t.Errorf("OrDefault() = %q, want %q", value, "default")
		}
	})

	t.Run("missing", func(t *testing.T) {
		const key = "TEST_OPTIONAL_MISSING"
		unsetEnv(t, key)

		if value := OrDefault(key, "default"); value != "default" {
			t.Errorf("OrDefault() = %q, want %q", value, "default")
		}
	})
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()

	previousValue, existed := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset environment variable %s: %v", key, err)
	}

	t.Cleanup(func() {
		if existed {
			_ = os.Setenv(key, previousValue)
			return
		}

		_ = os.Unsetenv(key)
	})
}
