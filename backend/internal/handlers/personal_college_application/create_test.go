package personalcollegeapplication

import (
	"context"
	"errors"
	"testing"
	"time"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the CreatePersonalCollegeApplication handler logic
// Written with stub student ID
func TestHandler_CreatePersonalCollegeApplication(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	edDeadline := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		input := models.CreatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: 1,
			ApplicationType: "ED",
			Category: "reach",
		}
		college := &models.GlobalCollege{
			ID: 1,
			SchoolName: "Northeastern University",
			EDDeadline: &edDeadline,
		}
		expectedApplication := &models.PersonalCollegeApplication{
			ID: 1,
			GlobalCollegeID: 1,
			ApplicationType: "ED",
			Category: "reach",
		}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPersonalRepo.On("CreatePersonalCollegeApplication", mock.Anything, mock.Anything, input).
			Return(expectedApplication, nil)

		handler := NewHandler(mockPersonalRepo, mockGlobalCollegeRepo)
		res, err := handler.CreatePersonalCollegeApplication(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expectedApplication, res)
	})

	t.Run("rejects deadline type the college does not offer", func(t *testing.T) {
		t.Parallel()

		input := models.CreatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: 1,
			ApplicationType: "EA",
			Category: "target",
		}
		// College only offers ED; EADeadline is nil.
		college := &models.GlobalCollege{
			ID: 1,
			SchoolName: "Northeastern University",
			EDDeadline: &edDeadline,
		}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		// CreatePersonalCollegeApplication is not called if the requested deadline is invalid

		handler := NewHandler(mockPersonalRepo, mockGlobalCollegeRepo)
		res, err := handler.CreatePersonalCollegeApplication(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("rejects unrecognized application_type", func(t *testing.T) {
		t.Parallel()

		input := models.CreatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: 1,
			ApplicationType: "NOT_REAL",
			Category: "reach",
		}
		college := &models.GlobalCollege{ID: 1, SchoolName: "Northeastern University"}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)

		handler := NewHandler(mockPersonalRepo, mockGlobalCollegeRepo)
		res, err := handler.CreatePersonalCollegeApplication(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("global college lookup error propagates", func(t *testing.T) {
		t.Parallel()

		input := models.CreatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: 999,
			ApplicationType: "ED",
			Category: "reach",
		}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(999)).
			Return(nil, errors.New("global college with id='999' not found"))

		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)

		handler := NewHandler(mockPersonalRepo, mockGlobalCollegeRepo)
		res, err := handler.CreatePersonalCollegeApplication(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "global college with id='999' not found")
	})

	t.Run("create repository error propagates", func(t *testing.T) {
		t.Parallel()

		input := models.CreatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: 1,
			ApplicationType: "ED",
			Category: "reach",
		}
		college := &models.GlobalCollege{ID: 1, EDDeadline: &edDeadline}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPersonalRepo.On("CreatePersonalCollegeApplication", mock.Anything, mock.Anything, input).
			Return(nil, errors.New("database error"))

		handler := NewHandler(mockPersonalRepo, mockGlobalCollegeRepo)
		res, err := handler.CreatePersonalCollegeApplication(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
