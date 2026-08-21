package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go-market/internal/platform/identity"

	"github.com/golang-jwt/jwt/v5"
)

const (
	minimumSecretLength = 32
	clockSkew           = 30 * time.Second
)

type tokenClaims struct {
	Roles []string `json:"roles,omitempty"`

	jwt.RegisteredClaims
}

type Verifier struct {
	secret   []byte
	issuer   string
	audience string
	now      func() time.Time
}

func NewVerifier(
	secret string,
	issuer string,
	audience string,
) (*Verifier, error) {
	if len(secret) < minimumSecretLength {
		return nil, fmt.Errorf(
			"JWT secret must contain at least %d bytes",
			minimumSecretLength,
		)
	}

	issuer = strings.TrimSpace(issuer)
	if issuer == "" {
		return nil, errors.New("JWT issuer is required")
	}

	audience = strings.TrimSpace(audience)
	if audience == "" {
		return nil, errors.New("JWT audience is required")
	}

	return &Verifier{
		secret:   []byte(secret),
		issuer:   issuer,
		audience: audience,
		now:      time.Now,
	}, nil
}

func (v *Verifier) Verify(rawToken string) (identity.Principal, error) {
	claims := &tokenClaims{}

	token, err := jwt.ParseWithClaims(
		rawToken,
		claims,
		func(*jwt.Token) (any, error) {
			return v.secret, nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(clockSkew),
		jwt.WithTimeFunc(v.now),
		jwt.WithStrictDecoding(),
	)
	if err != nil {
		return identity.Principal{}, fmt.Errorf("verify JWT: %w", err)
	}

	if !token.Valid {
		return identity.Principal{}, errors.New("JWT is invalid")
	}

	if strings.TrimSpace(claims.Subject) == "" {
		return identity.Principal{}, errors.New("JWT subject is required")
	}

	return identity.Principal{
		UserID: claims.Subject,
		Roles:  append([]string(nil), claims.Roles...),
	}, nil
}
