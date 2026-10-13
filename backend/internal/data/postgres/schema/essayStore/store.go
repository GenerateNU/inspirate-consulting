package essayRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EssayRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlEssayFiles embed.FS

func NewEssayRepository(db *pgxpool.Pool) *EssayRepository {
	return &EssayRepository{db: db}
}
