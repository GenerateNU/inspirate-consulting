// Isolated tests for AuthMiddleware: how each kind of missing, malformed,
// forged, expired, revoked or role-less credential is treated.
package routes

import (
	"encoding/base64"
	"encoding/json"
	"inspirate-consulting/internal/data"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// GET /user/{id} is guarded for both roles; with the repository mocked, an
// accepted request returns 200 and a rejected one never reaches the handler.
func newMiddlewareApp(t *testing.T) (*authApp, string) {
	t.Helper()
	user := &models.User{ID: uuid.New(), Name: "Target", SupabaseID: uuid.New()}
	repo := mocks.NewUserRepository(t)
	repo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: user.ID}).
		Return(&models.FetchUserOutput{Body: user}, nil).Maybe()
	return newAuthApp(t, func(r *data.Repository) { r.User = repo }), "/user/" + user.ID.String()
}

func TestAuthMiddleware_MissingCredentials(t *testing.T) {
	t.Parallel()

	a, path := newMiddlewareApp(t)
	valid := tokenFor(roleStudent)

	tests := []struct {
		name   string
		mutate func(*http.Request)
	}{
		{"no cookie", func(*http.Request) {}},
		{"empty jwt cookie", func(r *http.Request) { r.AddCookie(&http.Cookie{Name: "jwt", Value: ""}) }},
		{"token under a different cookie name", func(r *http.Request) {
			r.AddCookie(&http.Cookie{Name: "access_token", Value: valid})
		}},
		{"valid token in Authorization header instead of cookie", func(r *http.Request) {
			r.Header.Set("Authorization", "Bearer "+valid)
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := a.doReq(http.MethodGet, path, nil, tc.mutate)
			assert.Equal(t, http.StatusUnauthorized, resp.Status)
			assert.Equal(t, "Token Not Found", resp.Detail())
		})
	}
}

func TestAuthMiddleware_RejectsInvalidTokens(t *testing.T) {
	t.Parallel()

	a, path := newMiddlewareApp(t)
	sub := "22222222-2222-2222-2222-222222222222"
	attackerKey := newTestKey()

	expired := claimsFor(sub, roleCounselor)
	expired["iat"] = time.Now().Add(-2 * time.Hour).Unix()
	expired["exp"] = time.Now().Add(-1 * time.Hour).Unix()

	notYetValid := claimsFor(sub, roleCounselor)
	notYetValid["nbf"] = time.Now().Add(time.Hour).Unix()

	tests := []struct {
		name  string
		token string
	}{
		{"not a jwt", "not-a-jwt"},
		{"truncated jwt", sign(claimsFor(sub, roleCounselor))[:40]},
		{"signed by an unknown key with the trusted kid", signWith(jwt.SigningMethodRS256, attackerKey, testKeyID, claimsFor(sub, roleCounselor))},
		{"signed with an unknown kid", signWith(jwt.SigningMethodRS256, attackerKey, "other-kid", claimsFor(sub, roleCounselor))},
		{"HS256 instead of ES256/RS256", signWith(jwt.SigningMethodHS256, []byte("secret"), testKeyID, claimsFor(sub, roleCounselor))},
		{"unsigned alg none", signWith(jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, testKeyID, claimsFor(sub, roleCounselor))},
		{"expired", sign(expired)},
		{"not valid yet (nbf in the future)", sign(notYetValid)},
		{"payload tampered after signing", tamperRole(t, sign(claimsFor(sub, roleStudent)), roleCounselor)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := a.do(http.MethodGet, path, tc.token, nil)
			assert.Equal(t, http.StatusUnauthorized, resp.Status)
			assert.Equal(t, "Invalid/Expired Token", resp.Detail())
		})
	}
}

func TestAuthMiddleware_RejectsRevokedSession(t *testing.T) {
	t.Parallel()

	a, path := newMiddlewareApp(t)
	tok := tokenFor(roleStudent)
	require.Equal(t, http.StatusOK, a.do(http.MethodGet, path, tok, nil).Status)

	// signed out in Supabase; the JWT itself is still unexpired
	a.revoke(tok)
	resp := a.do(http.MethodGet, path, tok, nil)
	assert.Equal(t, http.StatusUnauthorized, resp.Status)
	assert.Equal(t, "Invalid/Expired Token", resp.Detail())
}

func TestAuthMiddleware_RequiresRoleClaim(t *testing.T) {
	t.Parallel()

	a, path := newMiddlewareApp(t)
	sub := "44444444-4444-4444-4444-444444444444"

	noRole := claimsFor(sub, "")

	// user_metadata is editable by the user, so it must never grant a role
	roleInUserMetadata := claimsFor(sub, "")
	roleInUserMetadata["user_metadata"] = map[string]any{"role": roleCounselor}

	nonStringRole := claimsFor(sub, "")
	nonStringRole["app_metadata"] = map[string]any{"role": []string{roleCounselor}}

	emptyRole := claimsFor(sub, "")
	emptyRole["app_metadata"] = map[string]any{"role": ""}

	tests := []struct {
		name   string
		claims jwt.MapClaims
	}{
		{"no role in app_metadata", noRole},
		{"role only in user-editable user_metadata", roleInUserMetadata},
		{"role is not a string", nonStringRole},
		{"empty role", emptyRole},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := a.do(http.MethodGet, path, sign(tc.claims), nil)
			assert.Equal(t, http.StatusUnauthorized, resp.Status)
			assert.Equal(t, "Token Missing Role", resp.Detail())
		})
	}
}

func TestAuthMiddleware_AcceptsValidSession(t *testing.T) {
	t.Parallel()

	a, path := newMiddlewareApp(t)
	for _, role := range []string{roleStudent, roleCounselor} {
		t.Run(role, func(t *testing.T) {
			resp := a.do(http.MethodGet, path, tokenFor(role), nil)
			assert.Equal(t, http.StatusOK, resp.Status, "body: %s", resp.Body)
		})
	}
}

// tamperRole rewrites app_metadata.role in a signed token's payload while
// keeping the original signature, as an attacker editing the cookie would.
func tamperRole(t *testing.T, token, role string) string {
	t.Helper()
	parts := strings.Split(token, ".")
	require.Len(t, parts, 3)
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)
	var claims map[string]any
	require.NoError(t, json.Unmarshal(raw, &claims))
	claims["app_metadata"] = map[string]any{"role": role}
	raw, err = json.Marshal(claims)
	require.NoError(t, err)
	parts[1] = base64.RawURLEncoding.EncodeToString(raw)
	return strings.Join(parts, ".")
}
