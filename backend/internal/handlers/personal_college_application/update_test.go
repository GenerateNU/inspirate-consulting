package personalcollegeapplication

import (
	"context"
	"errors"
	"testing"
	"time"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_UpdatePersonalCollegeApplication(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	ctx := auth.WithStudentID(context.Background(), studentID)
	applicationID := int64(42)
	deadline := time.Now().Add(24 * time.Hour)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		college := &models.GlobalCollege{
			ID:         1,
			SchoolName: "Test University",
			RDDeadline: &deadline,
		}
		input := models.UpdatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: 1,
			ApplicationType: "RD",
			Category:        "safety",
		}
		expected := &models.PersonalCollegeApplication{
			ID:              applicationID,
			GlobalCollegeID: 1,
			ApplicationType: "RD",
			Category:        "safety",
		}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("UpdatePersonalCollegeApplication", mock.Anything, studentID, applicationID, input).Return(expected, nil)

		handler := NewHandler(mockPCARepo, mockGlobalCollegeRepo)
		res, err := handler.UpdatePersonalCollegeApplication(ctx, applicationID, input)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("rejects application_type the college doesn't offer", func(t *testing.T) {
		t.Parallel()

		// College only offers RD : no ED/EA deadlines set.
		college := &models.GlobalCollege{
			ID:         1,
			SchoolName: "RD Only University",
			RDDeadline: &deadline,
		}
		input := models.UpdatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: 1,
			ApplicationType: "ED",
			Category:        "reach",
		}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		// UpdatePersonalCollegeApplication should NOT be called : validation should fail first.

		handler := NewHandler(mockPCARepo, mockGlobalCollegeRepo)
		res, err := handler.UpdatePersonalCollegeApplication(ctx, applicationID, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("global college not found", func(t *testing.T) {
		t.Parallel()

		input := models.UpdatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: 999,
			ApplicationType: "RD",
			Category:        "safety",
		}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(999)).
			Return(nil, errors.New("global college not found"))

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)

		handler := NewHandler(mockPCARepo, mockGlobalCollegeRepo)
		res, err := handler.UpdatePersonalCollegeApplication(ctx, applicationID, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("repository error propagates (e.g. application not found)", func(t *testing.T) {
		t.Parallel()

		college := &models.GlobalCollege{
			ID:         1,
			SchoolName: "Test University",
			RDDeadline: &deadline,
		}
		input := models.UpdatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: 1,
			ApplicationType: "RD",
			Category:        "safety",
		}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("UpdatePersonalCollegeApplication", mock.Anything, studentID, applicationID, input).
			Return(nil, errors.New("personal college application not found"))

		handler := NewHandler(mockPCARepo, mockGlobalCollegeRepo)
		res, err := handler.UpdatePersonalCollegeApplication(ctx, applicationID, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("missing student id returns error", func(t *testing.T) {
		t.Parallel()

		college := &models.GlobalCollege{
			ID:         1,
			SchoolName: "Test University",
			RDDeadline: &deadline,
		}
		input := models.UpdatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: 1,
			ApplicationType: "RD",
			Category:        "safety",
		}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		// UpdatePersonalCollegeApplication is not called without a student id

		handler := NewHandler(mockPCARepo, mockGlobalCollegeRepo)
		res, err := handler.UpdatePersonalCollegeApplication(context.Background(), applicationID, input)

		assert.ErrorIs(t, err, auth.ErrNoStudentID)
		assert.Nil(t, res)
	})
}
