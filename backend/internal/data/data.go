package data

import (
	"context"
	essayReviewRepository "inspirate-consulting/internal/data/postgres/schema/essayReviewStore"
	globalCollegeRepository "inspirate-consulting/internal/data/postgres/schema/globalCollegeStore"
	greetingRepository "inspirate-consulting/internal/data/postgres/schema/greetingStore"
	personalCollegeApplicationRepository "inspirate-consulting/internal/data/postgres/schema/personalCollegeApplicationStore"
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

// To represent the Personal College Application schema
type PersonalCollegeApplicationRepository interface {
	CreatePersonalCollegeApplication(ctx context.Context, studentID string, application models.CreatePersonalCollegeApplicationRequestBody) (*models.PersonalCollegeApplication, error)
	ListPersonalCollegeApplicationsByStudentID(ctx context.Context, studentID string) ([]models.PersonalCollegeApplication, error)
}

// To represent the Essay Review Transaction schema
type EssayReviewRepository interface {
	RequestEssayReview(ctx context.Context, input models.RequestEssayReviewRequestBody) (*models.EssayReviewTransaction, error)
	RefundEssayReview(ctx context.Context, id uuid.UUID) (*models.RefundEssayReviewResponseBody, error)
	CompleteEssayReview(ctx context.Context, id uuid.UUID) (*models.EssayReviewTransaction, error)
	GetEssayReviewStatus(ctx context.Context, essayID uuid.UUID) (*models.EssayReviewStatus, error)
	SetReviewBalance(ctx context.Context, id uuid.UUID, reviewBalance int) (*models.Student, error)
}

type Repository struct {
	db *pgxpool.Pool
	// For each interface, add a field here
	Greeting                   GreetingRepository
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
		GlobalCollege:              globalCollegeRepository.NewGlobalCollegeRepository(db),
		PersonalCollegeApplication: personalCollegeApplicationRepository.NewPersonalCollegeApplicationRepository(db),
		EssayReview:                essayReviewRepository.NewEssayReviewRepository(db),
	}
}
