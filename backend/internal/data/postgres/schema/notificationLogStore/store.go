package notificationLogRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationLogRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlNotificationLogFiles embed.FS

func NewNotificationLogRepository(db *pgxpool.Pool) *NotificationLogRepository {
	return &NotificationLogRepository{db: db}
}
