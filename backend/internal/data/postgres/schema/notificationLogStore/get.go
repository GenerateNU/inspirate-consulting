package notificationLogRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *NotificationLogRepository) GetPastDueTasks(ctx context.Context) ([]models.TaskNotification, error) {
	return r.getTaskNotifications(ctx, "get_past_due.sql")
}

func (r *NotificationLogRepository) GetWithinDueTasks(ctx context.Context) ([]models.TaskNotification, error) {
	return r.getTaskNotifications(ctx, "get_within_due.sql")
}

func (r *NotificationLogRepository) getTaskNotifications(ctx context.Context, sqlFile string) ([]models.TaskNotification, error) {
	selectQuery, err := schema.ReadSQLBaseScript(sqlFile, SqlNotificationLogFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.TaskNotification])
}
