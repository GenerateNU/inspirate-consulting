package mediaAccessRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MediaAccessRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlMediaAccessFiles embed.FS

func NewMediaAccessRepository(db *pgxpool.Pool) *MediaAccessRepository {
	return &MediaAccessRepository{db: db}
}
