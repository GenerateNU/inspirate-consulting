package auth

import (
	"errors"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
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

// RoleVerifier builds a route middleware that only lets the given roles through
type RoleVerifier func(api huma.API, roles ...models.Role) func(ctx huma.Context, next func(huma.Context))

func VerifyRole(api huma.API, roles ...models.Role) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

		role, _ := ctx.Context().Value("Role").(string)

		for _, r := range roles {
			if role == string(r) {
				next(ctx)
				return
			}
		}

		slog.Warn("role not permitted for route", "role", role, "path", ctx.Operation().Path)
		if err := huma.WriteErr(api, ctx, http.StatusForbidden, "Insufficient role"); err != nil {
			slog.Error("Failed to write error", "err", err)
		}
	}
}

func AuthMiddleware(api huma.API, verifier *Verifier, sb supabase.SupabaseInterface) func(ctx huma.Context, next func(huma.Context)) {
	skipPaths := map[string]bool{
		"/api/v1/health": true,
		"/user/login":    true,
	}

	return func(ctx huma.Context, next func(huma.Context)) {
		if skipPaths[ctx.Operation().Path] {
			next(ctx)
			return
		}

		accessToken := ""
		if cookie, err := huma.ReadCookie(ctx, AccessTokenCookieName); err == nil {
			accessToken = cookie.Value
		}

		refreshed := false
		if accessToken == "" {
			token, err := refreshSession(ctx, sb)
			if err != nil {
				err := huma.WriteErr(api, ctx, http.StatusUnauthorized, "Token Not Found")
				if err != nil {
					slog.Error("Failed to write error", "err", err)
				}
				return
			}
			accessToken, refreshed = token, true
		}

		claims, err := verifier.Verify(accessToken)

		// only an expired token is refreshed
		if err != nil && !refreshed && errors.Is(err, jwt.ErrTokenExpired) {
			if token, refreshErr := refreshSession(ctx, sb); refreshErr == nil {
				accessToken = token
				claims, err = verifier.Verify(accessToken)
			}
		}

		if err != nil {
			slog.Error("jwt verify failed", "err", err)
			err := huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid/Expired Token")
			if err != nil {
				slog.Error("Failed to write error", "err", err)
			}
			return
		}

		role, ok := claims.AppMetadata["role"].(string)
		if !ok || role == "" {
			slog.Error("Failed to parse role", "sub", claims.Sub)
			err := huma.WriteErr(api, ctx, http.StatusUnauthorized, "Token Missing Role")
			if err != nil {
				slog.Error("Failed to write error", "err", err)
			}
			return
		}

		if err := sb.SupabaseValidateSession(Client, accessToken); err != nil {
			slog.Error("session validation failed", "err", err)
			err := huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid/Expired Token")
			if err != nil {
				slog.Error("Failed to write error", "err", err)
			}
			return
		}

		ctx = huma.WithValue(ctx, "Supabase-ID", claims.Sub)
		ctx = huma.WithValue(ctx, "JWT", accessToken)
		ctx = huma.WithValue(ctx, "Role", role)

		next(ctx)
	}
}

func refreshSession(ctx huma.Context, sb supabase.SupabaseInterface) (string, error) {
	cookie, err := huma.ReadCookie(ctx, RefreshTokenCookieName)
	if err != nil || cookie.Value == "" {
		return "", ErrMissingToken
	}

	res, err := sb.SupabaseRefreshToken(Client, cookie.Value)
	if err != nil {
		slog.Error("token refresh failed", "err", err)
		expired := ExpiredRefreshTokenCookie()
		ctx.AppendHeader("Set-Cookie", expired.String())
		return "", err
	}
	if res.AccessToken == "" || res.RefreshToken == "" {
		return "", ErrInvalidToken
	}

	access := AccessTokenCookie(res.AccessToken, res.ExpiresIn)
	refresh := RefreshTokenCookie(res.RefreshToken)
	ctx.AppendHeader("Set-Cookie", access.String())
	ctx.AppendHeader("Set-Cookie", refresh.String())

	return res.AccessToken, nil
}
