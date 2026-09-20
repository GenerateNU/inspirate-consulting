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

// Unit tests for the GetGlobalCollege handler logic without HTTP or database dependencies.
func TestHandler_GetGlobalCollege(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	var id int64 = 1

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		expectedEntity := &models.GlobalCollege{
			ID: 1,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
			SchoolName: "Northeastern University",
			SchoolLocation: "Boston, MA",
		}

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("GetGlobalCollege", mock.Anything, id).Return(expectedEntity, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.GetGlobalCollege(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedEntity, res)
	})

	t.Run("not found propagates repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("GetGlobalCollege", mock.Anything, id).
			Return(nil, errors.New("global college with id='1' not found"))

		handler := NewHandler(mockRepo)
		res, err := handler.GetGlobalCollege(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "global college with id='1' not found")
	})
}