package data

import (
	"context"
	"time"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	essayGroupRepository "inspirate-consulting/internal/data/postgres/schema/essayGroupsStore"
	essayReviewRepository "inspirate-consulting/internal/data/postgres/schema/essayReviewStore"
	chatMessageRepository "inspirate-consulting/internal/data/postgres/schema/chatMessageStore"
	essayRepository "inspirate-consulting/internal/data/postgres/schema/essayStore"
	extracurricularRepository "inspirate-consulting/internal/data/postgres/schema/extracurricularStore"
	globalCollegeRepository "inspirate-consulting/internal/data/postgres/schema/globalCollegeStore"
	greetingRepository "inspirate-consulting/internal/data/postgres/schema/greetingStore"
	mediaAccessRepository "inspirate-consulting/internal/data/postgres/schema/mediaAccessStore"
	mediaRepository "inspirate-consulting/internal/data/postgres/schema/mediaStore"
	personalCollegeApplicationRepository "inspirate-consulting/internal/data/postgres/schema/personalCollegeApplicationStore"
	todoItemRepository "inspirate-consulting/internal/data/postgres/schema/todoItemStore"
	userRepository "inspirate-consulting/internal/data/postgres/schema/userStore"
	notificationPreferencesRepository "inspirate-consulting/internal/data/postgres/schema/notificationPreferencesStore"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// For each schema, their interfaces are to be defined here
type GreetingRepository interface {
	CreateGreeting(ctx context.Context, greeting models.CreateGreetingInput) (*models.CreateGreetingOutput, error)
}

type ExtracurricularRepository interface {
	CreateExtracurricular(ctx context.Context, userID string, extracurricular models.CreateExtracurricularInput) (*models.CreateExtracurricularOutput, error)
	ListExtracurriculars(ctx context.Context, studentID string) ([]models.Extracurricular, error)
	UpdateExtracurricular(ctx context.Context, id int64, extracurricular models.UpdateExtracurricularInput) (*models.Extracurricular, error)
}

// Essay Repository
type EssayRepository interface {
	UpdateStatus(ctx context.Context, essayID uuid.UUID, status models.Status) (*models.Essays, error)
	UpdateEssayGroup(ctx context.Context, essayID uuid.UUID, essayGroupID *uuid.UUID) (*models.Essays, error)
	GetEssaysFromStudent(ctx context.Context, studentID uuid.UUID) ([]models.Essays, error)
	GetEssaysByGroup(ctx context.Context, essayGroupID uuid.UUID) ([]models.Essays, error)
	CreateEssay(ctx context.Context, essay models.Essays) error
}

// Essay Group Repository
type EssayGroupRepository interface {
	CreateEssayGroup(ctx context.Context, group models.CreateEssayGroupBody) (*models.EssayGroups, error)
	ListEssayGroups(ctx context.Context, studentID uuid.UUID) ([]models.EssayGroups, error)
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
	FetchUser(ctx context.Context, input models.FetchUserInput) (*models.FetchUserOutput, error)
	FetchUserBySupabaseID(ctx context.Context, input models.FetchUserBySupabaseIDInput) (*models.FetchUserOutput, error)
	UpdateResetTime(ctx context.Context, id uuid.UUID, resetTime *time.Time) error
}

// To represent the methods for interacting with AWS S3 for presigned URLs for uploading and viewing videos
// Get and List generate fresh presigned URLs to avoid expiration issues
type VideoRepository interface {
	// generates a unique S3 key and a presigned PUT URL for a new video upload.
	PresignUpload(ctx context.Context, originalFilename string) (*models.PresignUploadResponse, error)
	GetVideo(ctx context.Context, s3Key string) (*models.Video, error)
	ListVideos(ctx context.Context) ([]models.Video, error)
}

type MediaRepository interface {
	CreateMedia(ctx context.Context, item *models.CreateMediaRequestBody) (*models.Media, error)
	GetMedia(ctx context.Context, id string) (*models.Media, error)
	ListAllMedia(ctx context.Context, limit int, offset int) ([]models.Media, error)
	DeleteMedia(ctx context.Context, id string) error
}

type MediaAccessRepository interface {
	GrantMediaAccess(ctx context.Context, body *models.GrantMediaAccessRequestBody) (*models.MediaAccess, error)
	RevokeMediaAccess(ctx context.Context, id string) error
	ListAccessibleMedia(ctx context.Context, studentID string, limit int, offset int) ([]models.Media, error)
}

// To represent the Essay Review Transaction schema
type EssayReviewRepository interface {
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

// To represent the Chat Message schema
type ChatMessageRepository interface {
	CreateChatMessage(ctx context.Context, senderID uuid.UUID, body models.CreateChatMessageRequestBody) (*models.ChatMessage, error)
	ListChats(ctx context.Context, userID uuid.UUID) ([]models.ChatSummary, error)
	ListChatMessages(ctx context.Context, userID uuid.UUID, otherUserID uuid.UUID, before *pagination.TimeIDKey, limit int) ([]models.ChatMessage, error)
	EditChatMessage(ctx context.Context, id uuid.UUID, senderID uuid.UUID, message string) (*models.ChatMessage, error)
	UpdateChatMessageReadAt(ctx context.Context, id uuid.UUID, recipientID uuid.UUID, readAt *time.Time) (*models.ChatMessage, error)
	MarkChatRead(ctx context.Context, userID uuid.UUID, otherUserID uuid.UUID) error
}

type NotificationPreferencesRepository interface {
	GetNotificationPreferences(ctx context.Context, userID uuid.UUID) (*models.NotificationPreferences, error)
	UpdateNotificationPreferences(ctx context.Context, userID uuid.UUID, preferences models.UpdateNotificationPreferencesRequestBody) (*models.NotificationPreferences, error)
}

type Repository struct {
	db *pgxpool.Pool

	// For each interface, add a field here
	Greeting                   GreetingRepository
	Extracurricular            ExtracurricularRepository
	Essay                      EssayRepository
	EssayGroup                 EssayGroupRepository
	ChatMessage                ChatMessageRepository
	TodoItem                   TodoItemRepository
	GlobalCollege              GlobalCollegeRepository
	PersonalCollegeApplication PersonalCollegeApplicationRepository
	Media                      MediaRepository
	MediaAccess                MediaAccessRepository
	User                       UserRepository
	Video                      VideoRepository
	EssayReview                EssayReviewRepository
	NotificationPreferences    NotificationPreferencesRepository
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
		db:                         db,
		Greeting:                   greetingRepository.NewGreetingRepository(db),
		Extracurricular:            extracurricularRepository.NewExtracurricularRepository(db),
		Essay:                      essayRepository.NewEssayRepository(db),
		EssayGroup:                 essayGroupRepository.NewEssayGroupRepository(db),
		ChatMessage:                chatMessageRepository.NewChatMessageRepository(db),
		TodoItem:                   todoItemRepository.NewTodoItemRepository(db),
		GlobalCollege:              globalCollegeRepository.NewGlobalCollegeRepository(db),
		PersonalCollegeApplication: personalCollegeApplicationRepository.NewPersonalCollegeApplicationRepository(db),
		Media:                      mediaRepository.NewMediaRepository(db),
		MediaAccess:                mediaAccessRepository.NewMediaAccessRepository(db),
		User:                       userRepository.NewUserRepository(db),
		EssayReview: essayReviewRepository.NewEssayReviewRepository(db),
		NotificationPreferences:    notificationPreferencesRepository.NewNotificationPreferencesRepository(db),
	}
}
