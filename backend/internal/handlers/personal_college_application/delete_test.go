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

func TestHandler_DeletePersonalCollegeApplication(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	ctx := auth.WithStudentID(context.Background(), studentID)
	applicationID := int64(42)
	collegeID := int64(7)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("DeletePersonalCollegeApplication", mock.Anything, studentID, applicationID).Return(collegeID, nil)

		mockActivityRepo := mocks.NewActivityRepository(t)
		mockActivityRepo.On("CreateActivity", mock.Anything, studentID, models.CreateActivityBody{
			Action:          models.ActionRemovedSchool,
			GlobalCollegeID: &collegeID,
		}).Return(nil)

		handler := NewHandler(mockPCARepo, nil, nil, mockActivityRepo)
		err := handler.DeletePersonalCollegeApplication(ctx, applicationID)

		assert.NoError(t, err)
	})

	t.Run("activity logging error is swallowed", func(t *testing.T) {
		t.Parallel()

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("DeletePersonalCollegeApplication", mock.Anything, studentID, applicationID).Return(collegeID, nil)

		mockActivityRepo := mocks.NewActivityRepository(t)
		mockActivityRepo.On("CreateActivity", mock.Anything, studentID, mock.Anything).
			Return(errors.New("activity write failed"))

		handler := NewHandler(mockPCARepo, nil, nil, mockActivityRepo)
		err := handler.DeletePersonalCollegeApplication(ctx, applicationID)

		assert.NoError(t, err)
	})

	t.Run("repository error propagates (e.g. not found)", func(t *testing.T) {
		t.Parallel()

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("DeletePersonalCollegeApplication", mock.Anything, studentID, applicationID).
			Return(int64(0), errors.New("personal college application not found"))

		handler := NewHandler(mockPCARepo, nil, nil, nil)
		err := handler.DeletePersonalCollegeApplication(ctx, applicationID)

		assert.Error(t, err)
		assert.EqualError(t, err, "personal college application not found")
	})

	t.Run("missing student id returns error", func(t *testing.T) {
		t.Parallel()

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		// DeletePersonalCollegeApplication is not called without a student id

		handler := NewHandler(mockPCARepo, nil, nil, nil)
		err := handler.DeletePersonalCollegeApplication(context.Background(), applicationID)

		assert.ErrorIs(t, err, auth.ErrNoStudentID)
	})
}
