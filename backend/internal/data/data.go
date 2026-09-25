package data

import (
	"context"
	globalCollegeRepository "inspirate-consulting/internal/data/postgres/schema/globalCollegeStore"
	essayRepository "inspirate-consulting/internal/data/postgres/schema/essayStore"
	greetingRepository "inspirate-consulting/internal/data/postgres/schema/greetingStore"
	personalCollegeApplicationRepository "inspirate-consulting/internal/data/postgres/schema/personalCollegeApplicationStore"
	mediaRepository "inspirate-consulting/internal/data/postgres/schema/mediaStore"
	mediaAccessRepository "inspirate-consulting/internal/data/postgres/schema/mediaAccessStore"
	todoItemRepository "inspirate-consulting/internal/data/postgres/schema/todoItemStore"
	userRepository "inspirate-consulting/internal/data/postgres/schema/userStore"
	"inspirate-consulting/internal/models"
	"time"

	"github.com/google/uuid"

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
type UserRepository interface {
	CreateUser(ctx context.Context, user models.CreateUserInput, supabase_id uuid.UUID) (*models.CreateUserOutput, error)
}

// To represent the methods for interacting with AWS S3 for presigned URLs for uploading and viewing videos
// Get and List generate fresh presigned URLs to avoid expiration issues
type VideoRepository interface {
	// generates a unique S3 key and a presigned PUT URL for a new video upload.
	PresignUpload(ctx context.Context, originalFilename string) (*models.PresignUploadResponse, error)
	GetVideo(ctx context.Context, s3Key string) (*models.Video, error)
	ListVideos(ctx context.Context) ([]models.Video, error)
}

type MediaRepository interface{
	CreateMedia(ctx context.Context, item *models.CreateMediaRequestBody)(*models.Media, error)
	GetMedia(ctx context.Context, id string)(*models.Media, error)
	ListAllMedia(ctx context.Context)([]models.Media, error)
	DeleteMedia(ctx context.Context, id string)error
}

type MediaAccessRepository interface{
	GrantMediaAccess(ctx context.Context, body *models.GrantMediaAccessRequestBody)(*models.MediaAccess,error)
	RevokeMediaAccess(ctx context.Context, id string)error
	ListAccessibleMedia(ctx context.Context, studentID string)([]models.Media, error)

// Essay Repository
type EssayRepository interface {
	UpdateStatus(ctx context.Context, essayID uuid.UUID, status models.Status) error
	GetEssaysFromStudent(ctx context.Context, studentID uuid.UUID) ([]models.Essays, error)
	CreateEssay(ctx context.Context, essay models.Essays) error
}

type Repository struct {
	db *pgxpool.Pool
	// For each interface, add a field here
	Greeting                   GreetingRepository
	GlobalCollege              GlobalCollegeRepository
	PersonalCollegeApplication PersonalCollegeApplicationRepository
	Media MediaRepository
	MediaAccess MediaAccessRepository
	TodoItem                   TodoItemRepository
	User                       UserRepository
	Video                      VideoRepository
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
		Greeting:                   greetingRepository.NewGreetingRepository(db),
		TodoItem:                   todoItemRepository.NewTodoItemRepository(db),
		GlobalCollege:              globalCollegeRepository.NewGlobalCollegeRepository(db),
		PersonalCollegeApplication: personalCollegeApplicationRepository.NewPersonalCollegeApplicationRepository(db),
		Media: mediaRepository.NewMediaRepository(db),
		MediaAccess: mediaAccessRepository.NewMediaAccessRepository(db),
		User:                       userRepository.NewUserRepository(db),
		Greeting: greetingRepository.NewGreetingRepository(db),
		Essay:    essayRepository.NewEssayRepository(db),
	}
}
