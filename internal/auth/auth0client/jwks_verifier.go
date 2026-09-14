package auth0client

import (
	"context"
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

// IDTokenClaims represents the required claims from the Auth0 ID Token.
type IDTokenClaims struct {
	jwt.RegisteredClaims
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

// IDTokenVerifier validates the Auth0 ID Token (RS256) using the tenant's JWKS,
// ensuring the issuer and audience match the Native App.
type IDTokenVerifier struct {
	jwks             keyfunc.Keyfunc
	expectedIssuer   string
	expectedAudience string
}

// NewIDTokenVerifier creates a verifier, automatically fetching and refreshing
// the JWKS in the background while ctx is alive.
func NewIDTokenVerifier(ctx context.Context, auth0Domain, nativeAppClientID string) (*IDTokenVerifier, error) {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", auth0Domain)

	jwks, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("gagal memuat JWKS Auth0: %w", err)
	}

	return &IDTokenVerifier{
		jwks:             jwks,
		expectedIssuer:   fmt.Sprintf("https://%s/", auth0Domain),
		expectedAudience: nativeAppClientID,
	}, nil
}

// Verify validates the token's signature, issuer, audience, and expiration.
func (v *IDTokenVerifier) Verify(rawIDToken string) (*IDTokenClaims, error) {
	claims := &IDTokenClaims{}

	token, err := jwt.ParseWithClaims(
		rawIDToken,
		claims,
		v.jwks.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(v.expectedIssuer),
		jwt.WithAudience(v.expectedAudience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("id token tidak valid: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("id token tidak valid")
	}

	return claims, nil
}
