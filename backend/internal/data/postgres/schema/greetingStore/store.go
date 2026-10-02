package greetingRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GreetingRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlGreetingFiles embed.FS

func NewGreetingRepository(db *pgxpool.Pool) *GreetingRepository {
	return &GreetingRepository{db: db}
}
