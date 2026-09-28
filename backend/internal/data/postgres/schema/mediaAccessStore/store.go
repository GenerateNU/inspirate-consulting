package mediaAccessRepository

import "github.com/jackc/pgx/v5/pgxpool"

type MediaAccessRepository struct {
	db *pgxpool.Pool
}

func NewMediaAccessRepository(db *pgxpool.Pool) *MediaAccessRepository{
	return &MediaAccessRepository{db: db}
}
