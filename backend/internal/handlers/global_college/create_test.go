package globalcollege

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

// Unit tests for the CreateGlobalCollege handler logic without HTTP or database dependencies.
func TestHandler_CreateGlobalCollege(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	edDeadline := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	input := &models.CreateGlobalCollegeInput{
		Body: models.CreateGlobalCollegeRequestBody{
			SchoolName: "Northeastern University",
			SchoolLocation: "Boston, MA",
			EDDeadline: &edDeadline,
		},
	}

	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEntity := &models.GlobalCollege{
		ID: 1,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		SchoolName: "Northeastern University",
		SchoolLocation: "Boston, MA",
		EDDeadline: &edDeadline,
	}
	expectedOutput := &models.CreateGlobalCollegeOutput{Body: *expectedEntity}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("CreateGlobalCollege", mock.Anything, input.Body).Return(expectedEntity, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.CreateGlobalCollege(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("repository error propagates (e.g. duplicate conflict from the DB)", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("CreateGlobalCollege", mock.Anything, input.Body).
			Return(nil, errors.New("global college with school_name/school_location='Northeastern University / Boston, MA' already exists"))

		handler := NewHandler(mockRepo)
		res, err := handler.CreateGlobalCollege(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "global college with school_name/school_location='Northeastern University / Boston, MA' already exists")
	})
}
