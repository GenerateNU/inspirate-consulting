package notificationPreferencesRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationPreferencesRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlNotificationPreferencesFiles embed.FS

func NewNotificationPreferencesRepository(db *pgxpool.Pool) *NotificationPreferencesRepository {
	return &NotificationPreferencesRepository{db: db}
}
