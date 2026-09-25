package data

import (
	"context"
	dbinterface "inspirate-consulting/internal/data/db-interface"
	essayReviewRepository "inspirate-consulting/internal/data/postgres/schema/essayReviewStore"
	globalCollegeRepository "inspirate-consulting/internal/data/postgres/schema/globalCollegeStore"
	greetingRepository "inspirate-consulting/internal/data/postgres/schema/greetingStore"
	personalCollegeApplicationRepository "inspirate-consulting/internal/data/postgres/schema/personalCollegeApplicationStore"
	todoItemRepository "inspirate-consulting/internal/data/postgres/schema/todoItemStore"
	"inspirate-consulting/internal/models"
	"time"

	"github.com/google/uuid"
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

// To represent the Personal College Application schema
type PersonalCollegeApplicationRepository interface {
	CreatePersonalCollegeApplication(ctx context.Context, studentID string, application models.CreatePersonalCollegeApplicationRequestBody) (*models.PersonalCollegeApplication, error)
	ListPersonalCollegeApplicationsByStudentID(ctx context.Context, studentID string) ([]models.PersonalCollegeApplication, error)
}

// To represent the Essay Review Transaction schema
type EssayReviewRepository interface {
	// WithTx runs fn inside one transaction; DB returns the pool for
	// single-statement reads. Together they let callers own the transaction
	// boundary without importing pgx.
	WithTx(ctx context.Context, fn func(db dbinterface.QueryInterface) error) error
	DB() dbinterface.QueryInterface

	LockStudent(ctx context.Context, db dbinterface.QueryInterface, id uuid.UUID) (*models.Student, error)
	SetStudentBalance(ctx context.Context, db dbinterface.QueryInterface, id uuid.UUID, reviewBalance int) error
	AdjustStudentBalance(ctx context.Context, db dbinterface.QueryInterface, id uuid.UUID, delta int) error

	LockTransaction(ctx context.Context, db dbinterface.QueryInterface, id uuid.UUID) (*models.EssayReviewTransaction, error)
	FindOpenReviewForEssay(ctx context.Context, db dbinterface.QueryInterface, essayID uuid.UUID) (*uuid.UUID, error)
	FindEssayReviewStatus(ctx context.Context, db dbinterface.QueryInterface, essayID uuid.UUID) (*models.EssayReviewStatus, error)
	InsertSpend(ctx context.Context, db dbinterface.QueryInterface, studentID, essayID uuid.UUID, amount int) (*models.EssayReviewTransaction, error)
	InsertRefund(ctx context.Context, db dbinterface.QueryInterface, charge *models.EssayReviewTransaction) (*models.EssayReviewTransaction, error)
	InsertAdjustment(ctx context.Context, db dbinterface.QueryInterface, studentID uuid.UUID, delta int) error
	LinkRefund(ctx context.Context, db dbinterface.QueryInterface, chargeID, reversalID uuid.UUID) (*models.EssayReviewTransaction, error)
	MarkCompleted(ctx context.Context, db dbinterface.QueryInterface, id uuid.UUID) (*models.EssayReviewTransaction, error)
}

type Repository struct {
	db *pgxpool.Pool
	// For each interface, add a field here
	Greeting                   GreetingRepository
	TodoItem                   TodoItemRepository
	GlobalCollege              GlobalCollegeRepository
	PersonalCollegeApplication PersonalCollegeApplicationRepository
	EssayReview                EssayReviewRepository
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
		Greeting:                   greetingRepository.NewGreetingRepository(db),
		TodoItem:                   todoItemRepository.NewTodoItemRepository(db),
		GlobalCollege:              globalCollegeRepository.NewGlobalCollegeRepository(db),
		PersonalCollegeApplication: personalCollegeApplicationRepository.NewPersonalCollegeApplicationRepository(db),
		EssayReview:                essayReviewRepository.NewEssayReviewRepository(db),
	}
}
