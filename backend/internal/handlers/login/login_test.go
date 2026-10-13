package login

import (
	"context"
	"errors"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
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

// supabaseLoggingIn returns a mock whose login succeeds for supabaseID.
func supabaseLoggingIn(supabaseID uuid.UUID) *supabase.MockSupabase {
	return &supabase.MockSupabase{
		LoginFunc: func(email, password string) (models.LoginResponse, error) {
			return models.LoginResponse{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
				User:         models.UserResponse{ID: supabaseID},
			}, nil
		},
	}
}

func TestHandler_Login(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success returns the users row and sets the session cookie", func(t *testing.T) {
		t.Parallel()
		supabaseID := uuid.New()
		user := &models.User{ID: uuid.New(), Name: "Sam", SupabaseID: supabaseID}

		repo := mocks.NewUserRepository(t)
		repo.On("FetchUserBySupabaseID", mock.Anything, models.FetchUserBySupabaseIDInput{SupabaseID: supabaseID}).
			Return(&models.FetchUserOutput{Body: user}, nil).Once()

		out, err := NewHandler(repo).Login(ctx, loginInput("sam@example.com", "Password1!"), supabaseLoggingIn(supabaseID))
		require.NoError(t, err)

		assert.Equal(t, user, out.Body)

		c := out.SetCookie
		assert.Equal(t, "jwt", c.Name)
		assert.Equal(t, "access-token", c.Value)
		assert.Equal(t, 3600, c.MaxAge)
		assert.Equal(t, "/", c.Path)
		assert.True(t, c.HttpOnly)
		assert.True(t, c.Secure)
		assert.Equal(t, http.SameSiteLaxMode, c.SameSite)
	})

	t.Run("pending password reset is passed through", func(t *testing.T) {
		t.Parallel()
		supabaseID := uuid.New()
		resetTime := time.Now().Add(-time.Second)
		repo := mocks.NewUserRepository(t)
		repo.On("FetchUserBySupabaseID", mock.Anything, mock.Anything).
			Return(&models.FetchUserOutput{Body: &models.User{SupabaseID: supabaseID, ResetTime: &resetTime}}, nil)

		out, err := NewHandler(repo).Login(ctx, loginInput("new@example.com", "Temp-Pass-1"), supabaseLoggingIn(supabaseID))
		require.NoError(t, err)
		require.NotNil(t, out.Body.ResetTime)
		assert.True(t, out.Body.ResetTime.Equal(resetTime))
	})

	t.Run("bad credentials propagate Supabase's error and skip the DB", func(t *testing.T) {
		t.Parallel()
		sb := &supabase.MockSupabase{
			LoginFunc: func(string, string) (models.LoginResponse, error) {
				return models.LoginResponse{}, errs.NewHTTPError(http.StatusBadRequest, errors.New("Invalid login credentials"))
			},
		}
		repo := mocks.NewUserRepository(t) // no expectations: must not be called

		out, err := NewHandler(repo).Login(ctx, loginInput("sam@example.com", "wrong"), sb)
		assert.Nil(t, out)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusBadRequest, httpErr.GetStatus())
	})

	t.Run("missing users row fails the login", func(t *testing.T) {
		t.Parallel()
		repo := mocks.NewUserRepository(t)
		repo.On("FetchUserBySupabaseID", mock.Anything, mock.Anything).
			Return(nil, errs.NotFound("User", "supabase_id", "x"))

		out, err := NewHandler(repo).Login(ctx, loginInput("orphan@example.com", "Password1!"), supabaseLoggingIn(uuid.New()))
		assert.Nil(t, out, "no cookie may be issued without a users row")
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.GetStatus())
	})
}
