package activityRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlActivityFiles embed.FS

func NewActivityRepository(db *pgxpool.Pool) *ActivityRepository {
	return &ActivityRepository{db: db}
}
