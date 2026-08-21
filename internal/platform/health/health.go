// Package health объединяет проверки готовности зависимостей приложения.
package health

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"google.golang.org/grpc"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

// Checker проверяет готовность одной или нескольких зависимостей.
type Checker interface {
	Check(ctx context.Context) error
}

// CheckFunc позволяет использовать функцию как Checker.
type CheckFunc func(ctx context.Context) error

// Check выполняет проверку готовности.
func (f CheckFunc) Check(ctx context.Context) error {
	return f(ctx)
}

// Group запускает именованные проверки параллельно.
type Group struct {
	checks map[string]Checker
}

// NewGroup создает группу проверок готовности.
func NewGroup(checks map[string]Checker) *Group {
	copyOfChecks := make(map[string]Checker, len(checks))
	for name, check := range checks {
		copyOfChecks[name] = check
	}

	return &Group{checks: copyOfChecks}
}

// Check выполняет все проверки и объединяет найденные ошибки.
func (g *Group) Check(ctx context.Context) error {
	type result struct {
		name string
		err  error
	}

	results := make(chan result, len(g.checks))

	for name, check := range g.checks {
		go func() {
			results <- result{
				name: name,
				err:  check.Check(ctx),
			}
		}()
	}

	failed := make([]result, 0, len(g.checks))
	for range g.checks {
		current := <-results
		if current.err != nil {
			failed = append(failed, current)
		}
	}

	sort.Slice(failed, func(i, j int) bool {
		return failed[i].name < failed[j].name
	})

	errorsByService := make([]error, 0, len(failed))
	for _, current := range failed {
		errorsByService = append(
			errorsByService,
			fmt.Errorf("%s: %w", current.name, current.err),
		)
	}

	return errors.Join(errorsByService...)
}

type grpcHealthClient interface {
	Check(
		ctx context.Context,
		request *healthv1.HealthCheckRequest,
		options ...grpc.CallOption,
	) (*healthv1.HealthCheckResponse, error)
}

// NewGRPCChecker создает проверку стандартного gRPC Health API.
func NewGRPCChecker(
	client grpcHealthClient,
	service string,
) Checker {
	return CheckFunc(func(ctx context.Context) error {
		response, err := client.Check(
			ctx,
			&healthv1.HealthCheckRequest{Service: service},
		)
		if err != nil {
			return fmt.Errorf("check gRPC health: %w", err)
		}
		if response.GetStatus() != healthv1.HealthCheckResponse_SERVING {
			return fmt.Errorf(
				"gRPC health status is %s",
				response.GetStatus(),
			)
		}

		return nil
	})
}
