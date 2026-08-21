package health_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	platformhealth "go-market/internal/platform/health"
)

func TestGroupReportsFailedChecks(t *testing.T) {
	group := platformhealth.NewGroup(map[string]platformhealth.Checker{
		"catalog": platformhealth.CheckFunc(func(context.Context) error {
			return nil
		}),
		"orders": platformhealth.CheckFunc(func(context.Context) error {
			return errors.New("not serving")
		}),
	})

	err := group.Check(context.Background())
	if err == nil {
		t.Fatal("Check() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "orders: not serving") {
		t.Errorf("Check() error = %q, want failed service name", err)
	}
}

func TestGroupSucceedsWhenAllChecksPass(t *testing.T) {
	group := platformhealth.NewGroup(map[string]platformhealth.Checker{
		"catalog": platformhealth.CheckFunc(func(context.Context) error {
			return nil
		}),
	})

	if err := group.Check(context.Background()); err != nil {
		t.Fatalf("Check() error = %v", err)
	}
}
