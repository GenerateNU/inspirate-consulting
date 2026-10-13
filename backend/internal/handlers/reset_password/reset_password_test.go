package resetpassword

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

// AuthMiddleware stores the caller's Supabase user id under this key.
func asUser(supabaseID string) context.Context {
	return context.WithValue(context.Background(), "Supabase-ID", supabaseID) //nolint:staticcheck // key set by AuthMiddleware
}

func resetInput(pw string) *models.ResetPasswordInput {
	in := &models.ResetPasswordInput{}
	in.Body.NewPassword = pw
	return in
}

func TestHandler_ResetPassword(t *testing.T) {
	t.Parallel()

	t.Run("updates the caller's own password and clears the pending reset", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		acct := sb.AddAccount("sam@example.com", "Temp-Pass-1", "student")

		repo := mocks.NewUserRepository(t)
		repo.On("UpdateResetTime", mock.Anything, acct.ID, (*time.Time)(nil)).Return(nil).Once()

		_, err := NewHandler(repo).ResetPassword(asUser(acct.ID.String()), resetInput("MyOwnPass1!"), sb)
		require.NoError(t, err)

		require.Len(t, sb.ResetPasswordCalls, 1)
		assert.Equal(t, sb.ResetPasswordCalls[0], authtest.ResetPasswordCall{UserID: acct.ID.String(), NewPassword: "MyOwnPass1!"})
		updated, _ := sb.Account("sam@example.com")
		assert.Equal(t, "MyOwnPass1!", updated.Password)
	})

	t.Run("no caller identity on the context", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		repo := mocks.NewUserRepository(t)

		_, err := NewHandler(repo).ResetPassword(context.Background(), resetInput("MyOwnPass1!"), sb)
		assert.Error(t, err)
		assert.Empty(t, sb.ResetPasswordCalls, "must not change any password without knowing who is asking")
	})

	t.Run("supabase failure leaves the pending reset in place", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		sb.ResetPasswordErr = errs.NewHTTPError(http.StatusUnprocessableEntity, errors.New("weak password"))
		repo := mocks.NewUserRepository(t) // UpdateResetTime must not be called

		_, err := NewHandler(repo).ResetPassword(asUser(uuid.NewString()), resetInput("x"), sb)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusUnprocessableEntity, httpErr.GetStatus())
	})

	t.Run("caller id that is not a uuid", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		sb.AddAccount("weird@example.com", "x", "student")
		repo := mocks.NewUserRepository(t)

		_, err := NewHandler(repo).ResetPassword(asUser("not-a-uuid"), resetInput("MyOwnPass1!"), sb)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.GetStatus(), "unknown account is rejected by Supabase first")
	})

	t.Run("database failure propagates", func(t *testing.T) {
		t.Parallel()
		sb := authtest.NewFakeSupabase(authtest.SharedSigner())
		acct := sb.AddAccount("sam@example.com", "Temp-Pass-1", "student")
		repo := mocks.NewUserRepository(t)
		repo.On("UpdateResetTime", mock.Anything, acct.ID, (*time.Time)(nil)).Return(errors.New("db down"))

		_, err := NewHandler(repo).ResetPassword(asUser(acct.ID.String()), resetInput("MyOwnPass1!"), sb)
		assert.EqualError(t, err, "db down")
	})
}

func TestHandler_ForcePasswordReset(t *testing.T) {
	t.Parallel()

	t.Run("flags the target account with a reset time in the past", func(t *testing.T) {
		t.Parallel()
		target := uuid.New()
		var got *time.Time
		repo := mocks.NewUserRepository(t)
		repo.On("UpdateResetTime", mock.Anything, target, mock.AnythingOfType("*time.Time")).
			Run(func(args mock.Arguments) { got = args.Get(2).(*time.Time) }).
			Return(nil).Once()

		before := time.Now()
		_, err := NewHandler(repo).ForcePasswordReset(context.Background(), &models.ForcePasswordResetInput{SupabaseID: target})
		require.NoError(t, err)

		require.NotNil(t, got)
		assert.True(t, got.Before(before), "reset time must already have passed so the next login sees it")
		assert.WithinDuration(t, before, *got, 5*time.Second)
	})

	t.Run("database failure propagates", func(t *testing.T) {
		t.Parallel()
		repo := mocks.NewUserRepository(t)
		repo.On("UpdateResetTime", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db down"))

		_, err := NewHandler(repo).ForcePasswordReset(context.Background(), &models.ForcePasswordResetInput{SupabaseID: uuid.New()})
		assert.EqualError(t, err, "db down")
	})
}
