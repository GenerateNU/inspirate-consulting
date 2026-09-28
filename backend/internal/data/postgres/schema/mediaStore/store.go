package mediaRepository

import "github.com/jackc/pgx/v5/pgxpool"

type MediaRepository struct {
	db *pgxpool.Pool
}

func NewMediaRepository(db *pgxpool.Pool) *MediaRepository{
	return &MediaRepository{db: db}
}