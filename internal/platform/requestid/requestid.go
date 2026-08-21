// Package requestid создает и переносит идентификатор запроса между транспортами.
package requestid

import (
	"context"

	"github.com/google/uuid"
)

const (
	// HTTPHeader содержит имя HTTP-заголовка с идентификатором запроса.
	HTTPHeader = "X-Request-ID"

	// GRPCMetadataKey содержит ключ gRPC metadata с идентификатором запроса.
	GRPCMetadataKey = "x-request-id"
)

type contextKey struct{}

// Normalize возвращает идентификатор запроса в UUID-формате.
// Если переданное значение некорректно, создаётся новый идентификатор.
func Normalize(value string) string {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.NewString()
	}

	return id.String()
}

// WithID добавляет идентификатор запроса в context.
func WithID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

// FromContext возвращает идентификатор запроса из context.
func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextKey{}).(string)
	return id, ok
}
