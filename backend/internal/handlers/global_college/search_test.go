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

	t.Run("success with results", func(t *testing.T) {
		t.Parallel()

		expected := []models.GlobalCollege{
			{ID: 1, SchoolName: "Northeastern University - Boston", SchoolLocation: "Boston, MA"},
		}

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("SearchGlobalColleges", mock.Anything, query, mock.Anything).Return(expected, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.SearchGlobalColleges(ctx, query)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("no results", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("SearchGlobalColleges", mock.Anything, query, mock.Anything).Return([]models.GlobalCollege{}, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.SearchGlobalColleges(ctx, query)

		assert.NoError(t, err)
		assert.Equal(t, []models.GlobalCollege{}, res)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("SearchGlobalColleges", mock.Anything, query, mock.Anything).
			Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.SearchGlobalColleges(ctx, query)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
