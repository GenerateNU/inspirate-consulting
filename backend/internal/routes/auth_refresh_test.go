// Tests for AuthMiddleware's refresh flow: a missing or expired access token is
// renewed from the refresh_token cookie, and the new cookies reach the browser.
package routes

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// addRefreshToken makes tok a refresh token Supabase accepts for a session of role.
func (a *authApp) addRefreshToken(tok, role string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.refreshTokens[tok] = role
}

// withSessionCookies sends the jwt and refresh_token cookies; an empty value is left out.
func withSessionCookies(accessToken, refreshToken string) func(*http.Request) {
	return func(r *http.Request) {
		if accessToken != "" {
			r.AddCookie(&http.Cookie{Name: "jwt", Value: accessToken})
		}
		if refreshToken != "" {
			r.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})
		}
	}
}

// lastCookie returns the last Set-Cookie for name, which is the one the browser keeps.
func lastCookie(resp apiResponse, name string) *http.Cookie {
	var found *http.Cookie
	for _, c := range resp.Cookies {
		if c.Name == name {
			found = c
		}
	}
	return found
}

func expiredToken(role string) string {
	claims := claimsFor(subFor(role), role)
	claims["iat"] = time.Now().Add(-2 * time.Hour).Unix()
	claims["exp"] = time.Now().Add(-1 * time.Hour).Unix()
	return sign(claims)
}

// assertNewSession checks the response set a fresh jwt and the rotated refresh token.
func assertNewSession(t *testing.T, resp apiResponse, oldRefresh string) {
	t.Helper()
	jwtCookie := lastCookie(resp, "jwt")
	require.NotNil(t, jwtCookie, "a new jwt cookie must be set")
	assert.NotEmpty(t, jwtCookie.Value)
	assert.Equal(t, 3600, jwtCookie.MaxAge)
	assert.Equal(t, "/", jwtCookie.Path)
	assert.True(t, jwtCookie.HttpOnly)

	refreshCookie := lastCookie(resp, "refresh_token")
	require.NotNil(t, refreshCookie, "the rotated refresh token must be set")
	assert.Equal(t, oldRefresh+"-rotated", refreshCookie.Value)
	assert.Positive(t, refreshCookie.MaxAge)
	assert.Equal(t, "/", refreshCookie.Path)
	assert.True(t, refreshCookie.HttpOnly)
}

func TestAuthMiddleware_RefreshesMissingAccessToken(t *testing.T) {
	t.Parallel()

	a, path := newMiddlewareApp(t)
	a.addRefreshToken("refresh-missing", roleStudent)

	resp := a.doReq(http.MethodGet, path, nil, withSessionCookies("", "refresh-missing"))
	assert.Equal(t, http.StatusOK, resp.Status, "body: %s", resp.Body)
	assertNewSession(t, resp, "refresh-missing")
}

func TestAuthMiddleware_RefreshesExpiredAccessToken(t *testing.T) {
	t.Parallel()

	a, path := newMiddlewareApp(t)
	a.addRefreshToken("refresh-expired", roleCounselor)

	resp := a.doReq(http.MethodGet, path, nil, withSessionCookies(expiredToken(roleCounselor), "refresh-expired"))
	assert.Equal(t, http.StatusOK, resp.Status, "body: %s", resp.Body)
	assertNewSession(t, resp, "refresh-expired")
}

func TestAuthMiddleware_RotatedRefreshTokenCannotBeReused(t *testing.T) {
	t.Parallel()

	a, path := newMiddlewareApp(t)
	a.addRefreshToken("refresh-once", roleStudent)

	require.Equal(t, http.StatusOK, a.doReq(http.MethodGet, path, nil, withSessionCookies("", "refresh-once")).Status)

	resp := a.doReq(http.MethodGet, path, nil, withSessionCookies("", "refresh-once"))
	assert.Equal(t, http.StatusUnauthorized, resp.Status)
}

func TestAuthMiddleware_RejectedRefreshTokenIsCleared(t *testing.T) {
	t.Parallel()

	a, path := newMiddlewareApp(t)

	tests := []struct {
		name        string
		accessToken string
		detail      string
	}{
		{"no access token", "", "Token Not Found"},
		{"expired access token", expiredToken(roleStudent), "Invalid/Expired Token"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := a.doReq(http.MethodGet, path, nil, withSessionCookies(tc.accessToken, "unknown-refresh-token"))
			assert.Equal(t, http.StatusUnauthorized, resp.Status)
			assert.Equal(t, tc.detail, resp.Detail())

			refreshCookie := lastCookie(resp, "refresh_token")
			require.NotNil(t, refreshCookie, "a rejected refresh token must be deleted")
			assert.Empty(t, refreshCookie.Value)
			assert.Negative(t, refreshCookie.MaxAge)
			assert.Nil(t, lastCookie(resp, "jwt"))
		})
	}
}

func TestAuthMiddleware_InvalidAccessTokenIsNotRefreshed(t *testing.T) {
	t.Parallel()

	a, path := newMiddlewareApp(t)
	a.addRefreshToken("refresh-valid", roleCounselor)

	// a forged token must not be traded for a real session, even with a good refresh token
	forged := tamperRole(t, sign(claimsFor(subFor(roleStudent), roleStudent)), roleCounselor)
	resp := a.doReq(http.MethodGet, path, nil, withSessionCookies(forged, "refresh-valid"))
	assert.Equal(t, http.StatusUnauthorized, resp.Status)
	assert.Equal(t, "Invalid/Expired Token", resp.Detail())
	assert.Zero(t, a.refreshCalls.Load())
	assert.Empty(t, resp.Cookies)
}

func TestAuthMiddleware_ValidAccessTokenIsNotRefreshed(t *testing.T) {
	t.Parallel()

	a, path := newMiddlewareApp(t)
	a.addRefreshToken("refresh-unused", roleStudent)

	resp := a.doReq(http.MethodGet, path, nil, withSessionCookies(tokenFor(roleStudent), "refresh-unused"))
	assert.Equal(t, http.StatusOK, resp.Status, "body: %s", resp.Body)
	assert.Zero(t, a.refreshCalls.Load())
	assert.Empty(t, resp.Cookies)
}

func TestLogout_ExpiresBothCookiesAfterRefresh(t *testing.T) {
	t.Parallel()

	a, _ := newMiddlewareApp(t)
	a.addRefreshToken("refresh-logout", roleStudent)

	// the access token has already expired, so the middleware refreshes before logging out
	resp := a.doReq(http.MethodPost, "/user/logout", nil, withSessionCookies("", "refresh-logout"))
	assert.Equal(t, http.StatusNoContent, resp.Status, "body: %s", resp.Body)

	for _, name := range []string{"jwt", "refresh_token"} {
		c := lastCookie(resp, name)
		require.NotNil(t, c, "%s must be expired on logout", name)
		assert.Empty(t, c.Value, name)
		assert.Negative(t, c.MaxAge, name)
	}
}
