package identity_test

import (
	"context"
	"testing"

	"go-market/internal/platform/identity"

	"google.golang.org/grpc/metadata"
)

func TestOutgoingIdentityCanBeReadAsIncoming(t *testing.T) {
	ctx := identity.AddToOutgoing(
		context.Background(),
		identity.Principal{
			UserID: "user-1",
			Roles:  []string{"buyer", "seller"},
		},
	)

	outgoing, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("outgoing metadata is missing")
	}

	incomingContext := metadata.NewIncomingContext(
		context.Background(),
		outgoing,
	)
	principal, ok := identity.FromIncoming(incomingContext)
	if !ok {
		t.Fatal("identity is missing from incoming metadata")
	}

	if principal.UserID != "user-1" {
		t.Errorf("UserID = %q, want %q", principal.UserID, "user-1")
	}
	if len(principal.Roles) != 2 {
		t.Errorf("role count = %d, want %d", len(principal.Roles), 2)
	}
}

func TestRemoveFromOutgoingRemovesUntrustedIdentity(t *testing.T) {
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(
			identity.UserIDMetadataKey,
			"forged-user",
			identity.RolesMetadataKey,
			"admin",
		),
	)

	ctx = identity.RemoveFromOutgoing(ctx)
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("outgoing metadata is missing")
	}

	if values := md.Get(identity.UserIDMetadataKey); len(values) != 0 {
		t.Errorf("user ID metadata = %v, want empty", values)
	}
	if values := md.Get(identity.RolesMetadataKey); len(values) != 0 {
		t.Errorf("roles metadata = %v, want empty", values)
	}
}
