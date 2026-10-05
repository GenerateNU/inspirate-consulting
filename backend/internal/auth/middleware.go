package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Std errors
var (
	ErrMissingToken  = errors.New("missing JWT token")
	ErrInvalidToken  = errors.New("invalid JWT token")
	ErrInvalidMethod = errors.New("unexpected JWT signing method")
)

// SupabaseClaims represents the JWT claims from Supabase Auth
type SupabaseClaims struct {
	Sub          string                 `json:"sub"`
	Email        string                 `json:"email"`
	Phone        string                 `json:"phone"`
	Role         string                 `json:"role"`
	Aud          string                 `json:"aud"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
	jwt.RegisteredClaims
}

// Verifier checks Supabase JWTs against the project's published public keys (JWKS)
type Verifier struct {
	jwks keyfunc.Keyfunc
}

// NewVerifier loads the signing keys from {supabaseURL}/auth/v1/.well-known/jwks.json.
// Keys are cached and refreshed in the background, so create this once at startup.
func NewVerifier(supabaseURL string) (*Verifier, error) {
	jwks, err := keyfunc.NewDefault([]string{supabaseURL + "/auth/v1/.well-known/jwks.json"})
	if err != nil {
		return nil, err
	}
	return &Verifier{jwks: jwks}, nil
}

// Verify validates a JWT token and returns the claims
func (v *Verifier) Verify(tokenString string) (*SupabaseClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&SupabaseClaims{},
		v.jwks.Keyfunc,
		jwt.WithValidMethods([]string{"ES256", "RS256"}),
	)

	if err != nil {
		return nil, errors.Join(ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*SupabaseClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func AuthMiddleware(api huma.API, verifier *Verifier) func(ctx huma.Context, next func(huma.Context)) {
	skipPaths := map[string]bool{
		"/api/v1/health": true,
		"/user/login":    true,
	}

	return func(ctx huma.Context, next func(huma.Context)) {
		if skipPaths[ctx.Operation().Path] {
			next(ctx)
			return
		}

		cookie, err := huma.ReadCookie(ctx, "jwt")
		if err != nil || cookie.Value == "" {
			err := huma.WriteErr(api, ctx, http.StatusUnauthorized, "Token Not Found")
			if err != nil {
				slog.Error("Failed to write error", "err", err)
			}
			return
		}

		claims, err := verifier.Verify(cookie.Value)

		if err != nil {
			slog.Error("jwt verify failed", "err", err)
			err := huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid/Expired Token")
			if err != nil {
				slog.Error("Failed to write error", "err", err)
			}
			return
		}

		//ctx.SetHeader("Supabase-ID", claims.Sub)
		ctx = huma.WithValue(ctx, "Supabase-ID", claims.Sub)

		// will be used for specifying the role of the user (student or counselor)
		//	ctx.SetHeader("Role", claims.AppMetadata["Role"].(string))
		next(ctx)
	}
}
