package data

import (
	"context"
	essayRepository "inspirate-consulting/internal/data/postgres/schema/essayStore"
	globalCollegeRepository "inspirate-consulting/internal/data/postgres/schema/globalCollegeStore"
	greetingRepository "inspirate-consulting/internal/data/postgres/schema/greetingStore"
	todoItemRepository "inspirate-consulting/internal/data/postgres/schema/todoItemStore"
	personalCollegeApplicationRepository "inspirate-consulting/internal/data/postgres/schema/personalCollegeApplicationStore"
	"inspirate-consulting/internal/models"
	"time"

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

// To represent the Personal College Application schema
type PersonalCollegeApplicationRepository interface {
	CreatePersonalCollegeApplication(ctx context.Context, studentID string, application models.CreatePersonalCollegeApplicationRequestBody) (*models.PersonalCollegeApplication, error)
	ListPersonalCollegeApplicationsByStudentID(ctx context.Context, studentID string) ([]models.PersonalCollegeApplication, error)
}

type Repository struct {
	db *pgxpool.Pool
	// For each interface, add a field here
	Greeting GreetingRepository
	Essay    EssayRepository
	TodoItem TodoItemRepository
	GlobalCollege GlobalCollegeRepository
	PersonalCollegeApplication PersonalCollegeApplicationRepository
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
		TodoItem: todoItemRepository.NewTodoItemRepository(db),
		GlobalCollege: globalCollegeRepository.NewGlobalCollegeRepository(db),
		PersonalCollegeApplication: personalCollegeApplicationRepository.NewPersonalCollegeApplicationRepository(db),
	}
}
