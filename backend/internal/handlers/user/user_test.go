package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
)

// mockUserRepository is a minimal mock for UserRepository.
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) CreateUser(ctx context.Context, user models.CreateUserInput, supabase_id uuid.UUID) (*models.CreateUserOutput, error) {
	args := m.Called(ctx, user, supabase_id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CreateUserOutput), args.Error(1)
}

func (m *mockUserRepository) FetchUser(ctx context.Context, input models.FetchUserInput) (*models.FetchUserOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FetchUserOutput), args.Error(1)
}

func TestHandler_CreateUser(t *testing.T) {
	t.Parallel()

	input := &models.CreateUserInput{}
	input.Body.Name = "Alice"
	input.Body.Email = "alice@example.com"
	input.Body.Password = "password123"

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := &mockUserRepository{}
		mockRepo.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.CreateUser(context.Background(), input, &supabase.MockSupabase{})

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
