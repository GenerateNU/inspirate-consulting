package essayGroupRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EssayGroupRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlEssayGroupFiles embed.FS

func NewEssayGroupRepository(db *pgxpool.Pool) *EssayGroupRepository {
	return &EssayGroupRepository{db: db}
}
