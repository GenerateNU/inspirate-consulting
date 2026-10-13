package user

import (
	"context"
	"errors"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the CreateGlobalCollege handler logic without HTTP or database dependencies.
func TestHandler_CreateGlobalCollege(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	test_uuid := uuid.MustParse("7bebfe6e-36ca-4343-a46c-d9430a951ba7")
	s1 := uuid.New()
	key := "pfp-key"

	expectedUser := models.User{
		ID:         test_uuid,
		Name:       "Aleng123",
		SupabaseID: s1,
		PfpKey:     &key,
	}

	input := models.CreateUserInput{}
	input.Body.Name = "Aleng123"
	input.Body.Email = "tt@gmail.com"
	input.Body.Password = "2976$$Alen$$"
	input.Body.PfpKey = &key

	expectedOutput := models.CreateUserOutput{Body: &expectedUser}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("CreateUser", mock.Anything, input).Return(&expectedOutput, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.CreateUser(ctx, &input, &supabase.MockSupabase{})

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, *res.Body)
	})

	t.Run("repository error propagates (e.g. duplicate conflict from the DB)", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("CreateUser", mock.Anything, input).
			Return(nil, errors.New("user already exists"))

		handler := NewHandler(mockRepo)
		res, err := handler.CreateUser(ctx, &input, &supabase.MockSupabase{})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "user already exists")
	})
}
