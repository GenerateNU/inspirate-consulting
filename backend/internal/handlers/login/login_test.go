package login

import (
	"context"
	"errors"
	"inspirate-consulting/internal/auth/authtest"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func loginInput(email, password string) *models.LoginInput {
	in := &models.LoginInput{}
	in.Body.Email = email
	in.Body.Password = password
	return in
}

func TestHandler_Login(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success returns the users row and sets the session cookie", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		acct := sb.AddAccount("sam@example.com", "Password1!", "student")
		user := &models.User{ID: uuid.New(), Name: "Sam", SupabaseID: acct.ID}

		repo := mocks.NewUserRepository(t)
		repo.On("FetchUserBySupabaseID", mock.Anything, models.FetchUserBySupabaseIDInput{SupabaseID: acct.ID}).
			Return(&models.FetchUserOutput{Body: user}, nil).Once()

		out, err := NewHandler(repo).Login(ctx, loginInput("sam@example.com", "Password1!"), sb)
		require.NoError(t, err)

		assert.Equal(t, user, out.Body)

		c := out.SetCookie
		assert.Equal(t, "jwt", c.Name)
		assert.NotEmpty(t, c.Value)
		require.NoError(t, sb.SupabaseValidateSession(nil, c.Value), "cookie must hold the Supabase access token")
		assert.Equal(t, 3600, c.MaxAge)
		assert.Equal(t, "/", c.Path)
		assert.True(t, c.HttpOnly)
		assert.True(t, c.Secure)
		assert.Equal(t, http.SameSiteLaxMode, c.SameSite)
	})

	t.Run("pending password reset is passed through", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		acct := sb.AddAccount("new@example.com", "Temp-Pass-1", "student")
		resetTime := mustTime(t)
		repo := mocks.NewUserRepository(t)
		repo.On("FetchUserBySupabaseID", mock.Anything, mock.Anything).
			Return(&models.FetchUserOutput{Body: &models.User{SupabaseID: acct.ID, ResetTime: &resetTime}}, nil)

		out, err := NewHandler(repo).Login(ctx, loginInput("new@example.com", "Temp-Pass-1"), sb)
		require.NoError(t, err)
		require.NotNil(t, out.Body.ResetTime)
		assert.True(t, out.Body.ResetTime.Equal(resetTime))
	})

	t.Run("bad credentials propagate Supabase's error and skip the DB", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		sb.AddAccount("sam@example.com", "Password1!", "student")
		repo := mocks.NewUserRepository(t)

		out, err := NewHandler(repo).Login(ctx, loginInput("sam@example.com", "wrong"), sb)
		assert.Nil(t, out)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusBadRequest, httpErr.GetStatus())
	})

	t.Run("supabase failure propagates", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		sb.LoginErr = errors.New("connection refused")

		out, err := NewHandler(mocks.NewUserRepository(t)).Login(ctx, loginInput("a@b.c", "x"), sb)
		assert.Nil(t, out)
		assert.EqualError(t, err, "connection refused")
	})

	t.Run("missing users row fails the login (no cookie returned)", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		sb.AddAccount("orphan@example.com", "Password1!", "student")
		repo := mocks.NewUserRepository(t)
		repo.On("FetchUserBySupabaseID", mock.Anything, mock.Anything).
			Return(nil, errs.NotFound("User", "supabase_id", "x"))

		out, err := NewHandler(repo).Login(ctx, loginInput("orphan@example.com", "Password1!"), sb)
		assert.Nil(t, out)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.GetStatus())
	})
}

func mustTime(t *testing.T) time.Time {
	t.Helper()
	return time.Now().Add(-time.Second).Truncate(time.Second)
}
