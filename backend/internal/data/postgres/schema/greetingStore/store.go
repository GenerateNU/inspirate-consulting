package greetingStore

import "github.com/jackc/pgx/v5/pgxpool"

type GreetingStore struct {
	db *pgxpool.Pool
}

func NewGreetingStore(db *pgxpool.Pool) *GreetingStore {
	return &GreetingStore{db: db}
}
