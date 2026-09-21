package data

import (
	"context"
	globalCollegeRepository "inspirate-consulting/internal/data/postgres/schema/globalCollegeStore"
	greetingRepository "inspirate-consulting/internal/data/postgres/schema/greetingStore"
	todoItemRepository "inspirate-consulting/internal/data/postgres/schema/todoItemStore"
	"inspirate-consulting/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// For each schema, their interfaces are to be defined here
type GreetingRepository interface {
	CreateGreeting(ctx context.Context, greeting models.CreateGreetingInput) (*models.CreateGreetingOutput, error)
}

// To represent Todo Item schema
type TodoItemRepository interface {
	CreateTodoItem(ctx context.Context, item *models.CreateTodoItemRequestBody) (*models.TodoItem, error)
	GetTodoItemsByStudent(ctx context.Context, studentID string) ([]models.TodoItem, error)
	UpdateTodoItemCompletedAt(ctx context.Context, id string, completedAt *time.Time) (*models.TodoItem, error)
}

// To represent the Global College schema
type GlobalCollegeRepository interface {
	CreateGlobalCollege(ctx context.Context, global_college models.CreateGlobalCollegeRequestBody) (*models.GlobalCollege, error)
	GetGlobalCollege(ctx context.Context, id int64) (*models.GlobalCollege, error)
	ListGlobalColleges(ctx context.Context) ([]models.GlobalCollege, error)
}

type Repository struct {
	db *pgxpool.Pool
	// For each interface, add a field here
	Greeting GreetingRepository
	TodoItem TodoItemRepository
	GlobalCollege GlobalCollegeRepository
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
		TodoItem: todoItemRepository.NewTodoItemRepository(db),
		GlobalCollege: globalCollegeRepository.NewGlobalCollegeRepository(db),
	}
}
