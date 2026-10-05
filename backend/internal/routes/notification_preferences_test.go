package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)


func setupTestAppWithNotificationPreferences(mockRepo data.NotificationPreferencesRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		NotificationPreferences: mockRepo,
	}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

func TestRoute_GetNotificationPreferences(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewNotificationPreferencesRepository(t)
		mockRepo.On("GetNotificationPreferences", mock.Anything, mock.Anything).Return(&models.NotificationPreferences{
			UserID: uuid.New(),
			EmailEnabled: true,
			WeeklySummaryEnabled: true,
			DueDateNotificationsEnabled: true,
			DaysBeforeDue: 1,
			NotifyPastDue: true,
		}, nil)

		app, err := setupTestAppWithNotificationPreferences(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/notification-preferences", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.NotificationPreferences
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.True(t, output.EmailEnabled)
		assert.Equal(t, 1, output.DaysBeforeDue)
	})
}

func TestRoute_UpdateNotificationPreferences(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewNotificationPreferencesRepository(t)
		mockRepo.On("UpdateNotificationPreferences", mock.Anything, mock.Anything, mock.MatchedBy(func(in models.UpdateNotificationPreferencesRequestBody) bool {
			return in.DaysBeforeDue == 5 && in.EmailEnabled == false
		})).Return(&models.NotificationPreferences{
			UserID: uuid.New(),
			EmailEnabled: false,
			WeeklySummaryEnabled: false,
			DueDateNotificationsEnabled: true,
			DaysBeforeDue: 5,
			NotifyPastDue: false,
		}, nil)

		app, err := setupTestAppWithNotificationPreferences(mockRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"email_enabled": false,
			"weekly_summary_enabled": false,
			"due_date_notifications_enabled": true,
			"days_before_due": 5,
			"notify_past_due": false,
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "/notification-preferences", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.NotificationPreferences
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Equal(t, 5, output.DaysBeforeDue)
		assert.False(t, output.EmailEnabled)
	})

	t.Run("validation error -> missing required field", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewNotificationPreferencesRepository(t)
		// Mock should not be called when Huma validation fails.

		app, err := setupTestAppWithNotificationPreferences(mockRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"email_enabled": true,
			// missing every other required field
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "/notification-preferences", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}