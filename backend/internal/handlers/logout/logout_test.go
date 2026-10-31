package logout

import (
	"context"
	"errors"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
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
		var revoked string
		sb := &supabase.MockSupabase{
			LogoutFunc: func(tok string) (models.LogoutResponse, error) {
				revoked = tok
				return models.LogoutResponse{}, nil
			},
		}

		out, err := NewHandler(mocks.NewUserRepository(t)).Logout(withJWT("session-token"), &models.LogoutInput{}, sb)
		require.NoError(t, err)
		assert.Equal(t, "session-token", revoked)

		c := out.SetCookie
		assert.Equal(t, "jwt", c.Name)
		assert.Empty(t, c.Value)
		assert.Equal(t, -1, c.MaxAge, "negative MaxAge makes Go send Max-Age=0 (delete); 0 would omit it")
		assert.Equal(t, "/", c.Path, "must match the login cookie's path")
		assert.True(t, c.HttpOnly)
		assert.True(t, c.Secure)
		assert.Equal(t, http.SameSiteLaxMode, c.SameSite)
	})

	for _, tc := range []struct {
		name string
		ctx  context.Context
	}{
		{"no token on the context", context.Background()},
		{"empty token", withJWT("")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			called := false
			sb := &supabase.MockSupabase{
				LogoutFunc: func(string) (models.LogoutResponse, error) {
					called = true
					return models.LogoutResponse{}, nil
				},
			}

			out, err := NewHandler(mocks.NewUserRepository(t)).Logout(tc.ctx, &models.LogoutInput{}, sb)
			assert.Nil(t, out)
			var httpErr errs.HTTPError
			require.ErrorAs(t, err, &httpErr)
			assert.Equal(t, http.StatusBadRequest, httpErr.GetStatus())
			assert.False(t, called, "Supabase must not be called without a token")
		})
	}

	t.Run("supabase failure propagates", func(t *testing.T) {
		t.Parallel()
		sb := &supabase.MockSupabase{
			LogoutFunc: func(string) (models.LogoutResponse, error) {
				return models.LogoutResponse{}, errors.New("supabase down")
			},
		}

		out, err := NewHandler(mocks.NewUserRepository(t)).Logout(withJWT("tok"), &models.LogoutInput{}, sb)
		assert.Nil(t, out)
		assert.EqualError(t, err, "supabase down")
	})
}
