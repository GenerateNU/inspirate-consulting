package storage

import (
	"context"

	"github.com/GenerateNU/inspirate-consulting/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GreetingStore interface {
	CreateGreeting(ctx context.Context, greeting models.CreateGreetingInput) (*models.CreateGreetingOutput, error)
}

type Store struct {
	db       *pgxpool.Pool
	Greeting GreetingStore
}
