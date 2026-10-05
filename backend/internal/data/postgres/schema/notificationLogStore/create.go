package notificationLogRepository

import (
	"context"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *NotificationLogRepository) CreateNotificationLog(ctx context.Context, input *models.CreateNotificationLogInput) (*models.NotificationLog, error) {
	createdLog := &models.NotificationLog{}

	insertQuery, err := schema.ReadSQLBaseScript("create_notification_log.sql", SqlNotificationLogFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		insertQuery,
		input.UserID,
		input.TaskID,
		input.NotificationType,
		input.SentAt,
	).Scan(
		&createdLog.ID,
		&createdLog.UserID,
		&createdLog.TaskID,
		&createdLog.NotificationType,
		&createdLog.SentAt,
		&createdLog.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return createdLog, nil
}
