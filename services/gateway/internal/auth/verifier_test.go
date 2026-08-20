package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	testSecret   = "0123456789abcdef0123456789abcdef"
	testIssuer   = "go-market-auth"
	testAudience = "go-market-api"
)

func TestNewVerifierRejectsShortSecret(t *testing.T) {
	_, err := NewVerifier(
		"short-secret",
		testIssuer,
		testAudience,
	)
	if err == nil {
		t.Fatal("expected error for short JWT secret")
	}
}

func TestVerifierVerify(t *testing.T) {
	now := time.Date(
		2026,
		time.August,
		20,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	validClaims := tokenClaims{
		Roles: []string{"customer"},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			Issuer:    testIssuer,
			Audience:  jwt.ClaimStrings{testAudience},
			IssuedAt:  jwt.NewNumericDate(now.Add(-time.Minute)),
			NotBefore: jwt.NewNumericDate(now.Add(-time.Minute)),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}

	expiredClaims := validClaims
	expiredClaims.ExpiresAt = jwt.NewNumericDate(
		now.Add(-time.Minute),
	)

	wrongIssuerClaims := validClaims
	wrongIssuerClaims.Issuer = "another-issuer"

	missingExpirationClaims := validClaims
	missingExpirationClaims.ExpiresAt = nil

	tests := []struct {
		name          string
		claims        tokenClaims
		signingMethod jwt.SigningMethod
		signingSecret string
		wantError     bool
	}{
		{
			name:          "valid token",
			claims:        validClaims,
			signingMethod: jwt.SigningMethodHS256,
			signingSecret: testSecret,
		},
		{
			name:          "expired token",
			claims:        expiredClaims,
			signingMethod: jwt.SigningMethodHS256,
			signingSecret: testSecret,
			wantError:     true,
		},
		{
			name:          "wrong issuer",
			claims:        wrongIssuerClaims,
			signingMethod: jwt.SigningMethodHS256,
			signingSecret: testSecret,
			wantError:     true,
		},
		{
			name:          "missing expiration",
			claims:        missingExpirationClaims,
			signingMethod: jwt.SigningMethodHS256,
			signingSecret: testSecret,
			wantError:     true,
		},
		{
			name:          "wrong signature",
			claims:        validClaims,
			signingMethod: jwt.SigningMethodHS256,
			signingSecret: "another-secret-with-at-least-32-bytes",
			wantError:     true,
		},
		{
			name:          "wrong algorithm",
			claims:        validClaims,
			signingMethod: jwt.SigningMethodHS384,
			signingSecret: testSecret,
			wantError:     true,
		},
	}

	verifier, err := NewVerifier(
		testSecret,
		testIssuer,
		testAudience,
	)
	if err != nil {
		t.Fatalf("NewVerifier() error = %v", err)
	}

	verifier.now = func() time.Time {
		return now
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rawToken := signToken(
				t,
				test.signingMethod,
				test.claims,
				test.signingSecret,
			)

			principal, err := verifier.Verify(rawToken)

			if test.wantError {
				if err == nil {
					t.Fatal("expected verification error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Verify() error = %v", err)
			}

			if principal.UserID != "user-123" {
				t.Errorf(
					"UserID = %q, want %q",
					principal.UserID,
					"user-123",
				)
			}
		})
	}
}

func signToken(
	t *testing.T,
	method jwt.SigningMethod,
	claims tokenClaims,
	secret string,
) string {
	t.Helper()

	token := jwt.NewWithClaims(method, claims)

	rawToken, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	return rawToken
}
