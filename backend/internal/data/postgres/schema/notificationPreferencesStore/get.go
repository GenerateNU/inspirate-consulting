package notificationPreferencesRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *NotificationPreferencesRepository) GetNotificationPreferences(ctx context.Context, userID uuid.UUID) (*models.NotificationPreferences, error) {
	getQuery, err := schema.ReadSQLBaseScript("get_notification_preferences_by_user.sql", SqlNotificationPreferencesFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, getQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notificationPreferences, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[models.NotificationPreferences])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// create a default row for this user if it does not exist
			return r.CreateDefaultNotificationPreferences(ctx, userID)
		}
		return nil, err
	}

	return &notificationPreferences, nil
}
