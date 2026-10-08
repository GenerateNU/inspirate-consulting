package notificationpreferences

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_UpdateNotificationPreferences(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := models.UpdateNotificationPreferencesRequestBody{
		EmailEnabled: false,
		WeeklySummaryEnabled: false,
		DueDateNotificationsEnabled: true,
		DaysBeforeDue: 5,
		NotifyPastDue: false,
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expected := &models.NotificationPreferences{
			UserID: uuid.New(),
			EmailEnabled: false,
			WeeklySummaryEnabled:false,
			DueDateNotificationsEnabled: true,
			DaysBeforeDue: 5,
			NotifyPastDue: false,
		}

		mockRepo := mocks.NewNotificationPreferencesRepository(t)
		mockRepo.On("UpdateNotificationPreferences", mock.Anything, mock.Anything, input).Return(expected, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.UpdateNotificationPreferences(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewNotificationPreferencesRepository(t)
		mockRepo.On("UpdateNotificationPreferences", mock.Anything, mock.Anything, input).
			Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.UpdateNotificationPreferences(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}