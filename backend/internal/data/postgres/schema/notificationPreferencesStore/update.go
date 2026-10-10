package notificationPreferencesRepository

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *NotificationPreferencesRepository) UpdateNotificationPreferences(ctx context.Context, userID uuid.UUID, input models.UpdateNotificationPreferencesRequestBody) (*models.NotificationPreferences, error) {
	updatedPreferences := &models.NotificationPreferences{}

	updateQuery, err := schema.ReadSQLBaseScript("update_notification_preferences_by_user.sql", SqlNotificationPreferencesFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		updateQuery,
		userID,
		input.EmailEnabled,
		input.WeeklySummaryEnabled,
		input.DueDateNotificationsEnabled,
		input.DaysBeforeDue,
		input.NotifyPastDue,
	).Scan(
		&updatedPreferences.UserID,
		&updatedPreferences.EmailEnabled,
		&updatedPreferences.WeeklySummaryEnabled,
		&updatedPreferences.DueDateNotificationsEnabled,
		&updatedPreferences.DaysBeforeDue,
		&updatedPreferences.NotifyPastDue,
	)
	if err != nil {
		return nil, err
	}

	return updatedPreferences, nil
}
