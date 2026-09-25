package data

import (
	"context"
	essayRepository "inspirate-consulting/internal/data/postgres/schema/essayStore"
	greetingRepository "inspirate-consulting/internal/data/postgres/schema/greetingStore"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// For each schema, their interfaces are to be defined here
type GreetingRepository interface {
	CreateGreeting(ctx context.Context, greeting models.CreateGreetingInput) (*models.CreateGreetingOutput, error)
}

// Essay Repository
type EssayRepository interface {
	UpdateStatus(ctx context.Context, essayID uuid.UUID, status models.Status) error
	GetEssaysFromStudent(ctx context.Context, studentID uuid.UUID) ([]models.Essays, error)
	CreateEssay(ctx context.Context, essay models.Essays) error
}

type Repository struct {
	db *pgxpool.Pool
	// For each interface, add a field here
	Greeting GreetingRepository
	Essay    EssayRepository
}

// Close closes the database connection pool
func (r *Repository) Close() error {
	r.db.Close()
	return nil
}

// GetDB returns the underlying pgxpool.Pool instance
func (r *Repository) GetDB() *pgxpool.Pool {
	return r.db
}

// NewRepository creates a new Repository instance with the given database pool
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
		// For each interface, add an instance of the interface here
		Greeting: greetingRepository.NewGreetingRepository(db),
		Essay:    essayRepository.NewEssayRepository(db),
	}
}
