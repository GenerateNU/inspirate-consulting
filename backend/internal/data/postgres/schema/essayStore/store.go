package essayRepository

import "github.com/jackc/pgx/v5/pgxpool"

type EssayRepository struct {
	db *pgxpool.Pool
}

func NewEssayRepository(db *pgxpool.Pool) *EssayRepository {
	return &EssayRepository{db: db}
}
