package personalcollegeapplication

import (
	"context"
	"errors"
	"testing"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ListPersonalCollegeApplicationsByStudentID(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	ctx := auth.WithStudentID(context.Background(), studentID)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expectedApplications := []models.PersonalCollegeApplication{
			{ID: 1, GlobalCollegeID: 1, ApplicationType: "ED", Category: "reach"},
			{ID: 2, GlobalCollegeID: 2, ApplicationType: "EA", Category: "target"},
		}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPersonalRepo.On("ListPersonalCollegeApplicationsByStudentID", mock.Anything, studentID).
			Return(expectedApplications, nil)

		handler := NewHandler(mockPersonalRepo, mockGlobalCollegeRepo)
		res, err := handler.ListPersonalCollegeApplicationsByStudentID(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedApplications, res)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPersonalRepo.On("ListPersonalCollegeApplicationsByStudentID", mock.Anything, studentID).
			Return([]models.PersonalCollegeApplication{}, nil)

		handler := NewHandler(mockPersonalRepo, mockGlobalCollegeRepo)
		res, err := handler.ListPersonalCollegeApplicationsByStudentID(ctx)

		assert.NoError(t, err)
		assert.Equal(t, []models.PersonalCollegeApplication{}, res)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPersonalRepo.On("ListPersonalCollegeApplicationsByStudentID", mock.Anything, studentID).
			Return(nil, errors.New("database error"))

		handler := NewHandler(mockPersonalRepo, mockGlobalCollegeRepo)
		res, err := handler.ListPersonalCollegeApplicationsByStudentID(ctx)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})

	t.Run("missing student id returns error", func(t *testing.T) {
		t.Parallel()

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		// ListPersonalCollegeApplicationsByStudentID is not called without a student id

		handler := NewHandler(mockPersonalRepo, mockGlobalCollegeRepo)
		res, err := handler.ListPersonalCollegeApplicationsByStudentID(context.Background())

		assert.ErrorIs(t, err, auth.ErrNoStudentID)
		assert.Nil(t, res)
	})
}
