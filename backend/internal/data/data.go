package data

import (
	"context"
	globalCollegeRepository "inspirate-consulting/internal/data/postgres/schema/globalCollegeStore"
	greetingRepository "inspirate-consulting/internal/data/postgres/schema/greetingStore"
	userRepository "inspirate-consulting/internal/data/postgres/schema/userStore"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

// For each schema, their interfaces are to be defined here
type GreetingRepository interface {
	CreateGreeting(ctx context.Context, greeting models.CreateGreetingInput) (*models.CreateGreetingOutput, error)
}

// To represent the Global College schema
type GlobalCollegeRepository interface {
	CreateGlobalCollege(ctx context.Context, global_college models.CreateGlobalCollegeRequestBody) (*models.GlobalCollege, error)
	GetGlobalCollege(ctx context.Context, id int64) (*models.GlobalCollege, error)
	ListGlobalColleges(ctx context.Context) ([]models.GlobalCollege, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user models.CreateUserInput, supabase_id uuid.UUID) (*models.CreateUserOutput, error)
	FetchUser(ctx context.Context, user models.FetchUserInput) (*models.FetchUserOutput, error)
}

type Repository struct {
	db *pgxpool.Pool
	// For each interface, add a field here
	Greeting      GreetingRepository
	GlobalCollege GlobalCollegeRepository
	User          UserRepository
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
		Greeting:      greetingRepository.NewGreetingRepository(db),
		GlobalCollege: globalCollegeRepository.NewGlobalCollegeRepository(db),
		User:          userRepository.NewUserRepository(db),
	}
}
