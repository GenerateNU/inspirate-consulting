package user

import (
	"context"
	"errors"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_FetchUserBySupabaseID(t *testing.T) {
	t.Parallel()

	test_uuid := uuid.MustParse("7bebfe6e-36ca-4343-a46c-d9430a951ba7")
	s1 := uuid.New()
	key := "pfp-key"

	expectedUser := models.User{
		ID:         test_uuid,
		Name:       "Aleng123",
		SupabaseID: s1,
		PfpKey:     &key,
	}

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		input := models.FetchUserBySupabaseIDInput{SupabaseID: s1}
		expectedOutput := models.FetchUserOutput{Body: &expectedUser}

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUserBySupabaseID", mock.Anything, input).Return(&expectedOutput, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.FetchUserBySupabaseID(ctx, &input)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, *res.Body)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		input := models.FetchUserBySupabaseIDInput{SupabaseID: s1}

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUserBySupabaseID", mock.Anything, input).Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.FetchUserBySupabaseID(ctx, &input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
