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
	"github.com/stretchr/testify/require"
)

func TestHandler_FetchUser(t *testing.T) {
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

		input := models.FetchUserInput{ID: test_uuid}
		expectedOutput := models.FetchUserOutput{Body: &expectedUser}

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, input).Return(&expectedOutput, nil)
		pfpRepo := mocks.NewPfpRepository(t)
		pfpRepo.On("GetDownloadURL", ctx, key).Return("https://bucket/a.jpg", nil)

		handler := NewHandler(mockRepo, nil, pfpRepo)
		res, err := handler.FetchUser(ctx, &input)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, res.Body.User)
		require.NotNil(t, res.Body.PfpURL)
		assert.Equal(t, "https://bucket/a.jpg", *res.Body.PfpURL)
	})

	t.Run("user without a picture has no url", func(t *testing.T) {
		t.Parallel()

		input := models.FetchUserInput{ID: test_uuid}
		noPfp := models.User{ID: test_uuid, Name: "Aleng123"}

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, input).
			Return(&models.FetchUserOutput{Body: &noPfp}, nil)
		pfpRepo := mocks.NewPfpRepository(t)

		res, err := NewHandler(mockRepo, nil, pfpRepo).FetchUser(ctx, &input)

		require.NoError(t, err)
		assert.Equal(t, noPfp, res.Body.User)
		assert.Nil(t, res.Body.PfpURL)
	})

	t.Run("failed presign still returns the user", func(t *testing.T) {
		t.Parallel()

		input := models.FetchUserInput{ID: test_uuid}
		expectedOutput := models.FetchUserOutput{Body: &expectedUser}

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, input).Return(&expectedOutput, nil)
		pfpRepo := mocks.NewPfpRepository(t)
		pfpRepo.On("GetDownloadURL", ctx, key).Return("", errors.New("s3 unavailable"))

		res, err := NewHandler(mockRepo, nil, pfpRepo).FetchUser(ctx, &input)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, res.Body.User)
		assert.Nil(t, res.Body.PfpURL)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		input := models.FetchUserInput{ID: test_uuid}

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, input).Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo, nil, nil)
		res, err := handler.FetchUser(ctx, &input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
