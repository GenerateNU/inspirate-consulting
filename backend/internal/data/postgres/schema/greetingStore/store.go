package greetingRepository

import "github.com/jackc/pgx/v5/pgxpool"

type GreetingRepository struct {
	db *pgxpool.Pool
}

func NewGreetingRepository(db *pgxpool.Pool) *GreetingRepository {
	return &GreetingRepository{db: db}
}
