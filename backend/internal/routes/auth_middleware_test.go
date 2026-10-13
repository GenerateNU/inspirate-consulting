// Isolated tests for AuthMiddleware: how each kind of missing, malformed,
// forged, expired, revoked or role-less credential is treated.
package routes

import (
	"encoding/base64"
	"encoding/json"
	"inspirate-consulting/internal/auth/authtest"
	"inspirate-consulting/internal/models"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GET /user/{id} is guarded for both roles and backed by the in-memory repo,
// so an accepted request returns 200 with the user and a rejected one never
// reaches the handler.
func targetUser(t *testing.T, a *authApp) (*authtest.Account, string) {
	t.Helper()
	acct := a.seedUser("Target", "target@example.com", "Password1!", "student")
	u, ok := a.users.BySupabaseID(acct.ID)
	require.True(t, ok)
	return acct, "/user/" + u.ID.String()
}

func TestAuthMiddleware_MissingCredentials(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	_, path := targetUser(t, a)
	valid := a.sb.IssueSession("11111111-1111-1111-1111-111111111111", "student")

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

	a := newAuthApp(t, nil)
	_, path := targetUser(t, a)
	sub := "22222222-2222-2222-2222-222222222222"
	good := authtest.ClaimsFor(sub, "counselor")

	expired := good
	expired.IssuedAt = time.Now().Add(-2 * time.Hour)
	expired.ExpiresAt = time.Now().Add(-1 * time.Hour)

	notYetValid := good
	notYetValid.NotBefore = time.Now().Add(time.Hour)

	attackerKey := authtest.NewKey()

	tests := []struct {
		name  string
		token string
	}{
		{"not a jwt", "not-a-jwt"},
		{"truncated jwt", a.sb.Signer.Sign(good)[:40]},
		{"signed by an unknown key with the trusted kid", authtest.SignWith(attackerKey, authtest.KeyID, good, 1)},
		{"signed with an unknown kid", authtest.SignWith(attackerKey, "other-kid", good, 1)},
		{"HS256 instead of ES256/RS256", authtest.SignHS256([]byte("secret"), good)},
		{"unsigned alg none", authtest.SignNone(good)},
		{"expired", a.sb.Signer.Sign(expired)},
		{"not valid yet (nbf in the future)", a.sb.Signer.Sign(notYetValid)},
		{"payload tampered after signing", tamperRole(t, a.sb.Signer.Sign(authtest.ClaimsFor(sub, "student")), "counselor")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// mark the token live on the Supabase side so only the JWT checks can reject it
			a.sb.RegisterSession(tc.token)
			resp := a.do(http.MethodGet, path, tc.token, nil)
			assert.Equal(t, http.StatusUnauthorized, resp.Status)
			assert.Equal(t, "Invalid/Expired Token", resp.Detail())
		})
	}
}

func TestAuthMiddleware_RejectsDeadSessions(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	_, path := targetUser(t, a)
	sub := "33333333-3333-3333-3333-333333333333"

	t.Run("validly signed token Supabase has no session for", func(t *testing.T) {
		tok := a.sb.Signer.Sign(authtest.ClaimsFor(sub, "student"))
		resp := a.do(http.MethodGet, path, tok, nil)
		assert.Equal(t, http.StatusUnauthorized, resp.Status)
		assert.Equal(t, "Invalid/Expired Token", resp.Detail())
	})

	t.Run("session revoked on the Supabase side", func(t *testing.T) {
		tok := a.sb.IssueSession(sub, "student")
		require.Equal(t, http.StatusOK, a.do(http.MethodGet, path, tok, nil).Status)

		a.sb.Revoke(tok)
		resp := a.do(http.MethodGet, path, tok, nil)
		assert.Equal(t, http.StatusUnauthorized, resp.Status)
		assert.Equal(t, "Invalid/Expired Token", resp.Detail())
	})
}

func TestAuthMiddleware_RequiresRoleClaim(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	_, path := targetUser(t, a)
	sub := "44444444-4444-4444-4444-444444444444"

	noRole := authtest.ClaimsFor(sub, "")

	// user_metadata is editable by the user, so it must never grant a role
	roleInUserMetadata := authtest.ClaimsFor(sub, "")
	roleInUserMetadata.UserMetadata = map[string]any{"role": "counselor"}

	nonStringRole := authtest.ClaimsFor(sub, "")
	nonStringRole.AppMetadata["role"] = []string{"counselor"}

	emptyRole := authtest.ClaimsFor(sub, "")
	emptyRole.AppMetadata["role"] = ""

	tests := []struct {
		name   string
		claims authtest.Claims
	}{
		{"no role in app_metadata", noRole},
		{"role only in user-editable user_metadata", roleInUserMetadata},
		{"role is not a string", nonStringRole},
		{"empty role", emptyRole},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tok := a.sb.IssueSessionWithClaims(tc.claims)
			resp := a.do(http.MethodGet, path, tok, nil)
			assert.Equal(t, http.StatusUnauthorized, resp.Status)
			assert.Equal(t, "Token Missing Role", resp.Detail())
		})
	}
}

func TestAuthMiddleware_AcceptsValidSession(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	acct, path := targetUser(t, a)

	for _, role := range []string{"student", "counselor"} {
		t.Run(role, func(t *testing.T) {
			tok := a.sb.IssueSession(acct.ID.String(), role)
			resp := a.do(http.MethodGet, path, tok, nil)
			require.Equal(t, http.StatusOK, resp.Status, "body: %s", resp.Body)

			var u models.User
			resp.JSON(t, &u)
			assert.Equal(t, acct.ID, u.SupabaseID)
		})
	}
}

func TestAuthMiddleware_PublicRoutes(t *testing.T) {
	t.Parallel()

	a := newAuthApp(t, nil)
	a.seedUser("Pub", "pub@example.com", "Password1!", "student")

	t.Run("health needs no session", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, a.do(http.MethodGet, "/api/v1/health", "", nil).Status)
	})

	t.Run("login needs no session", func(t *testing.T) {
		resp, tok := a.login("pub@example.com", "Password1!")
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.NotEmpty(t, tok)
	})

	t.Run("login ignores a stale or garbage jwt cookie", func(t *testing.T) {
		resp := a.do(http.MethodPost, "/user/login", "garbage",
			map[string]string{"email": "pub@example.com", "password": "Password1!"})
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.NotNil(t, resp.Cookie("jwt"))
	})
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
