package notificationPreferencesRepository

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *NotificationPreferencesRepository) CreateDefaultNotificationPreferences(ctx context.Context, userID uuid.UUID) (*models.NotificationPreferences, error) {
	createdPreferences := &models.NotificationPreferences{}

	createQuery, err := schema.ReadSQLBaseScript("create_default_notification_preferences_by_user.sql", SqlNotificationPreferencesFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, createQuery, userID).Scan(
		&createdPreferences.UserID,
		&createdPreferences.EmailEnabled,
		&createdPreferences.WeeklySummaryEnabled,
		&createdPreferences.DueDateNotificationsEnabled,
		&createdPreferences.DaysBeforeDue,
		&createdPreferences.NotifyPastDue,
	)
	if err != nil {
		return nil, err
	}

	return createdPreferences, nil
}
