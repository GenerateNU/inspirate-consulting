package user

import (
	"context"
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

func TestHandler_FetchCurrentUser(t *testing.T) {
	t.Parallel()

	t.Run("counselor gets their counselor id", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		counselorID := uuid.New()
		ctx := auth.WithCounselorID(auth.WithUserID(context.Background(), userID), counselorID)

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
			Return(&models.FetchUserOutput{Body: &models.User{ID: userID, Name: "Test Counselor"}}, nil)

		res, err := NewHandler(mockRepo).FetchCurrentUser(ctx)

		require.NoError(t, err)
		assert.Equal(t, userID, res.ID)
		require.NotNil(t, res.CounselorID)
		assert.Equal(t, counselorID, *res.CounselorID)
		assert.Nil(t, res.StudentID)
	})

	t.Run("student gets their student id", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		studentID := uuid.New()
		ctx := auth.WithStudentID(auth.WithUserID(context.Background(), userID), studentID)

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
			Return(&models.FetchUserOutput{Body: &models.User{ID: userID, Name: "Test Student"}}, nil)

		res, err := NewHandler(mockRepo).FetchCurrentUser(ctx)

		require.NoError(t, err)
		require.NotNil(t, res.StudentID)
		assert.Equal(t, studentID, *res.StudentID)
		assert.Nil(t, res.CounselorID)
	})

	t.Run("no user in context is unauthorized", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewUserRepository(t)

		res, err := NewHandler(mockRepo).FetchCurrentUser(context.Background())

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 401, err.(errs.HTTPError).Code)
		mockRepo.AssertNotCalled(t, "FetchUser", mock.Anything, mock.Anything)
	})
}
