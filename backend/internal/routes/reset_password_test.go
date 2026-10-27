package routes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRoute_SendTempPassword(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	supabaseID := uuid.New()

	user := &models.User{
		ID:         userID,
		Name:       "Ada Lovelace",
		Email:      "ada@example.com",
		SupabaseID: supabaseID,
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
			Return(&models.FetchUserOutput{Body: user}, nil)
		mockRepo.On("UpdateResetTime", mock.Anything, supabaseID, mock.Anything).Return(nil)

		app, err := setupTestAppWithUser(mockRepo, &supabase.MockSupabase{})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/user/"+userID.String()+"/send-temp-password", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.SendTempPasswordBody
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Equal(t, "ada@example.com", output.Email)
	})

	t.Run("user not found", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
			Return(nil, errors.New("user not found"))

		app, err := setupTestAppWithUser(mockRepo, &supabase.MockSupabase{})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/user/"+userID.String()+"/send-temp-password", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
