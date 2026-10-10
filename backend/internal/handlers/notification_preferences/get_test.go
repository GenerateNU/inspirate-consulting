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

func TestHandler_GetNotificationPreferences(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expected := &models.NotificationPreferences{
			UserID: uuid.New(),
			EmailEnabled: true,
			WeeklySummaryEnabled: true,
			DueDateNotificationsEnabled: true,
			DaysBeforeDue: 1,
			NotifyPastDue: true,
		}

		mockRepo := mocks.NewNotificationPreferencesRepository(t)
		mockRepo.On("GetNotificationPreferences", mock.Anything, mock.Anything).Return(expected, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.GetNotificationPreferences(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewNotificationPreferencesRepository(t)
		mockRepo.On("GetNotificationPreferences", mock.Anything, mock.Anything).
			Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.GetNotificationPreferences(ctx)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}