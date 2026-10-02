package mediaRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MediaRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlMediaFiles embed.FS

func NewMediaRepository(db *pgxpool.Pool) *MediaRepository {
	return &MediaRepository{db: db}
}
