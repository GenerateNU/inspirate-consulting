package resetpassword

import (
	"context"
	"errors"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
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
		caller := uuid.New()
		var gotUser, gotPassword string
		sb := &supabase.MockSupabase{
			ResetPasswordFunc: func(newPassword, userID string) (models.ResetPasswordResponse, error) {
				gotPassword, gotUser = newPassword, userID
				return models.ResetPasswordResponse{}, nil
			},
		}
		repo := mocks.NewUserRepository(t)
		repo.On("UpdateResetTime", mock.Anything, caller, (*time.Time)(nil)).Return(nil).Once()

		_, err := NewHandler(repo).ResetPassword(asUser(caller.String()), resetInput("MyOwnPass1!"), sb)
		require.NoError(t, err)
		assert.Equal(t, caller.String(), gotUser, "only the caller's own account may be changed")
		assert.Equal(t, "MyOwnPass1!", gotPassword)
	})

	t.Run("no caller identity on the context", func(t *testing.T) {
		t.Parallel()
		called := false
		sb := &supabase.MockSupabase{
			ResetPasswordFunc: func(string, string) (models.ResetPasswordResponse, error) {
				called = true
				return models.ResetPasswordResponse{}, nil
			},
		}

		_, err := NewHandler(mocks.NewUserRepository(t)).ResetPassword(context.Background(), resetInput("MyOwnPass1!"), sb)
		assert.Error(t, err)
		assert.False(t, called, "must not change any password without knowing who is asking")
	})

	t.Run("supabase failure leaves the pending reset in place", func(t *testing.T) {
		t.Parallel()
		sb := &supabase.MockSupabase{
			ResetPasswordFunc: func(string, string) (models.ResetPasswordResponse, error) {
				return models.ResetPasswordResponse{}, errors.New("weak password")
			},
		}
		repo := mocks.NewUserRepository(t) // UpdateResetTime must not be called

		_, err := NewHandler(repo).ResetPassword(asUser(uuid.NewString()), resetInput("x"), sb)
		assert.EqualError(t, err, "weak password")
	})

	t.Run("database failure propagates", func(t *testing.T) {
		t.Parallel()
		repo := mocks.NewUserRepository(t)
		repo.On("UpdateResetTime", mock.Anything, mock.Anything, (*time.Time)(nil)).Return(errors.New("db down"))

		_, err := NewHandler(repo).ResetPassword(asUser(uuid.NewString()), resetInput("MyOwnPass1!"), &supabase.MockSupabase{})
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
