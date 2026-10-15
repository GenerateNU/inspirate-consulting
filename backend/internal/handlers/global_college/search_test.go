package globalcollege

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_SearchGlobalColleges(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	query := "Northeastern"

	t.Run("passes the search threshold and result limit to the repository", func(t *testing.T) {
		t.Parallel()

		expected := []models.GlobalCollege{
			{ID: 1, SchoolName: "Northeastern University - Boston", SchoolLocation: "Boston, MA"},
		}

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("SearchGlobalColleges", mock.Anything, query, searchThreshold, searchResultLimit).Return(expected, nil)

		res, err := NewHandler(mockRepo).SearchGlobalColleges(ctx, query)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("no results", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("SearchGlobalColleges", mock.Anything, query, searchThreshold, searchResultLimit).
			Return([]models.GlobalCollege{}, nil)

		res, err := NewHandler(mockRepo).SearchGlobalColleges(ctx, query)

		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("SearchGlobalColleges", mock.Anything, query, searchThreshold, searchResultLimit).
			Return(nil, errors.New("database error"))

		res, err := NewHandler(mockRepo).SearchGlobalColleges(ctx, query)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}