package auth_test

import (
	"testing"

	gatewayauth "go-market/gateway/internal/auth"
	catalogv1 "go-market/gen/go/catalog/v1"

	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

func TestPolicyIsPublic(t *testing.T) {
	tests := []struct {
		name       string
		fullMethod string
		want       bool
		wantError  bool
	}{
		{
			name:       "annotated public method",
			fullMethod: catalogv1.CatalogService_ListProducts_FullMethodName,
			want:       true,
		},
		{
			name:       "method without annotation is private",
			fullMethod: healthv1.Health_Check_FullMethodName,
			want:       false,
		},
		{
			name:       "unknown method",
			fullMethod: "/unknown.v1.Unknown/Get",
			wantError:  true,
		},
		{
			name:       "invalid method name",
			fullMethod: "invalid",
			wantError:  true,
		},
	}

	policy := gatewayauth.Policy{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := policy.IsPublic(test.fullMethod)

			if test.wantError {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("IsPublic() error = %v", err)
			}

			if got != test.want {
				t.Errorf("IsPublic() = %v, want %v", got, test.want)
			}
		})
	}
}
