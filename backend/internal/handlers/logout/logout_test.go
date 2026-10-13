package logout

import (
	"context"
	"errors"
	"inspirate-consulting/internal/auth/authtest"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AuthMiddleware stores the session token under this key.
func withJWT(token string) context.Context {
	return context.WithValue(context.Background(), "JWT", token) //nolint:staticcheck // key set by AuthMiddleware
}

func TestHandler_Logout(t *testing.T) {
	t.Parallel()

	t.Run("revokes the caller's session and expires the cookie", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		tok := sb.IssueSession("11111111-1111-1111-1111-111111111111", "student")

		out, err := NewHandler(mocks.NewUserRepository(t)).Logout(withJWT(tok), &models.LogoutInput{}, sb)
		require.NoError(t, err)

		assert.Equal(t, []string{tok}, sb.LogoutCalls)
		assert.Error(t, sb.SupabaseValidateSession(nil, tok), "session must be dead after logout")

		c := out.SetCookie
		assert.Equal(t, "jwt", c.Name)
		assert.Empty(t, c.Value)
		assert.Equal(t, -1, c.MaxAge, "negative MaxAge makes Go send Max-Age=0 (delete); 0 would omit it")
		assert.Equal(t, "/", c.Path, "must match the login cookie's path")
		assert.True(t, c.HttpOnly)
		assert.True(t, c.Secure)
		assert.Equal(t, http.SameSiteLaxMode, c.SameSite)
		assert.Contains(t, c.String(), "Max-Age=0")
	})

	for _, tc := range []struct {
		name string
		ctx  context.Context
	}{
		{"no token on the context", context.Background()},
		{"empty token", withJWT("")},
		{"token of the wrong type", context.WithValue(context.Background(), "JWT", 42)}, //nolint:staticcheck
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sb := authtest.NewFakeSupabase(authtest.SharedSigner())

			out, err := NewHandler(mocks.NewUserRepository(t)).Logout(tc.ctx, &models.LogoutInput{}, sb)
			assert.Nil(t, out)
			var httpErr errs.HTTPError
			require.ErrorAs(t, err, &httpErr)
			assert.Equal(t, http.StatusBadRequest, httpErr.GetStatus())
			assert.Empty(t, sb.LogoutCalls, "Supabase must not be called without a token")
		})
	}

	t.Run("supabase failure propagates and no cookie is cleared", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		sb.LogoutErr = errors.New("supabase down")

		out, err := NewHandler(mocks.NewUserRepository(t)).Logout(withJWT("tok"), &models.LogoutInput{}, sb)
		assert.Nil(t, out)
		assert.EqualError(t, err, "supabase down")
	})
}
