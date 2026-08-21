// Package identity переносит проверенную identity между внутренними gRPC-сервисами.
package identity

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	// UserIDMetadataKey содержит идентификатор пользователя во внутренней metadata.
	UserIDMetadataKey = "x-user-id"

	// RolesMetadataKey содержит роли пользователя во внутренней metadata.
	RolesMetadataKey = "x-user-roles"
)

// Principal описывает проверенную identity пользователя.
type Principal struct {
	UserID string
	Roles  []string
}

// RemoveFromOutgoing удаляет identity, которую мог передать внешний клиент.
func RemoveFromOutgoing(ctx context.Context) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ctx
	}

	md = md.Copy()
	md.Delete(UserIDMetadataKey)
	md.Delete(RolesMetadataKey)

	return metadata.NewOutgoingContext(ctx, md)
}

// AddToOutgoing добавляет проверенную identity в исходящую gRPC metadata.
func AddToOutgoing(
	ctx context.Context,
	principal Principal,
) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	} else {
		md = md.Copy()
	}

	md.Set(UserIDMetadataKey, principal.UserID)
	if len(principal.Roles) > 0 {
		md.Set(RolesMetadataKey, principal.Roles...)
	}

	return metadata.NewOutgoingContext(ctx, md)
}

// FromIncoming возвращает identity из входящей gRPC metadata.
// Вызывать эту функцию можно только для запросов из доверенной внутренней сети.
func FromIncoming(ctx context.Context) (Principal, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return Principal{}, false
	}

	userIDs := md.Get(UserIDMetadataKey)
	if len(userIDs) != 1 || strings.TrimSpace(userIDs[0]) == "" {
		return Principal{}, false
	}

	return Principal{
		UserID: userIDs[0],
		Roles:  append([]string(nil), md.Get(RolesMetadataKey)...),
	}, true
}
