package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func TestHandler_GetCurrentUser(t *testing.T) {
	t.Parallel()

	t.Run("counselor gets their counselor id", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		counselorID := uuid.New()
		ctx := auth.WithCounselorID(auth.WithUserID(context.Background(), userID), counselorID)

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
			Return(&models.FetchUserOutput{Body: &models.User{ID: userID, Name: "Test Counselor"}}, nil)

		res, err := NewHandler(mockRepo, nil, nil).GetCurrentUser(ctx)

		require.NoError(t, err)
		assert.Equal(t, userID, res.Body.ID)
		assert.Equal(t, models.CounselorRole, res.Body.Role)
		require.NotNil(t, res.Body.CounselorID)
		assert.Equal(t, counselorID, *res.Body.CounselorID)
		assert.Nil(t, res.Body.StudentID)
	})

	t.Run("student gets their student id", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		studentID := uuid.New()
		ctx := auth.WithStudentID(auth.WithUserID(context.Background(), userID), studentID)

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
			Return(&models.FetchUserOutput{Body: &models.User{ID: userID, Name: "Test Student"}}, nil)

		res, err := NewHandler(mockRepo, nil, nil).GetCurrentUser(ctx)

		require.NoError(t, err)
		assert.Equal(t, models.StudentRole, res.Body.Role)
		require.NotNil(t, res.Body.StudentID)
		assert.Equal(t, studentID, *res.Body.StudentID)
		assert.Nil(t, res.Body.CounselorID)
	})

	t.Run("no user in context is unauthorized", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewUserRepository(t)

		res, err := NewHandler(mockRepo, nil, nil).GetCurrentUser(context.Background())

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 401, err.(errs.HTTPError).Code)
		mockRepo.AssertNotCalled(t, "FetchUser", mock.Anything, mock.Anything)
	})

	t.Run("user with a picture gets a download url", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		ctx := auth.WithUserID(context.Background(), userID)
		key := "pfps/" + userID.String() + "/a.jpg"

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
			Return(&models.FetchUserOutput{Body: &models.User{ID: userID, PfpKey: &key}}, nil)
		pfpRepo := mocks.NewPfpRepository(t)
		pfpRepo.On("GetDownloadURL", ctx, key).Return("https://bucket/a.jpg", nil)

		res, err := NewHandler(mockRepo, nil, pfpRepo).GetCurrentUser(ctx)

		require.NoError(t, err)
		require.NotNil(t, res.Body.PfpURL)
		assert.Equal(t, "https://bucket/a.jpg", *res.Body.PfpURL)
	})

	t.Run("failed presign still returns the user", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		ctx := auth.WithUserID(context.Background(), userID)
		key := "pfps/" + userID.String() + "/a.jpg"

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
			Return(&models.FetchUserOutput{Body: &models.User{ID: userID, PfpKey: &key}}, nil)
		pfpRepo := mocks.NewPfpRepository(t)
		pfpRepo.On("GetDownloadURL", ctx, key).Return("", errors.New("s3 unavailable"))

		res, err := NewHandler(mockRepo, nil, pfpRepo).GetCurrentUser(ctx)

		require.NoError(t, err)
		assert.Equal(t, userID, res.Body.ID)
		assert.Nil(t, res.Body.PfpURL)
	})
}
