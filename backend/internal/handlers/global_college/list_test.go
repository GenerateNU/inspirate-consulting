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

// Unit tests for the ListGlobalColleges handler logic without HTTP or database dependencies.
func TestHandler_ListGlobalColleges(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := &models.ListGlobalCollegesInput{}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		expectedColleges := []models.GlobalCollege{
			{ID: 1, CreatedAt: createdAt, UpdatedAt: createdAt, SchoolName: "Alpha University", SchoolLocation: "Alpha City, AA"},
			{ID: 2, CreatedAt: createdAt, UpdatedAt: createdAt, SchoolName: "Beta College", SchoolLocation: "Beta Town, BB"},
		}
		expectedOutput := &models.ListGlobalCollegesOutput{Body: expectedColleges}

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("ListGlobalColleges", mock.Anything).Return(expectedColleges, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.ListGlobalColleges(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("ListGlobalColleges", mock.Anything).Return([]models.GlobalCollege{}, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.ListGlobalColleges(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, &models.ListGlobalCollegesOutput{Body: []models.GlobalCollege{}}, res)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("ListGlobalColleges", mock.Anything).Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.ListGlobalColleges(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}