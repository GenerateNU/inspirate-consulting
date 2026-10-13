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

func TestHandler_UpdateGlobalCollege(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	collegeID := int64(1)
	input := models.UpdateGlobalCollegeRequestBody{
		SchoolName:     "Updated University",
		SchoolLocation: "Updated City, UC",
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expected := &models.GlobalCollege{
			ID:             collegeID,
			SchoolName:     "Updated University",
			SchoolLocation: "Updated City, UC",
		}

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("UpdateGlobalCollege", mock.Anything, collegeID, input, mock.Anything).Return(expected, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.UpdateGlobalCollege(ctx, collegeID, input)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("repository error propagates (e.g. near-duplicate conflict or not found)", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("UpdateGlobalCollege", mock.Anything, collegeID, input, mock.Anything).
			Return(nil, errors.New("global college not found"))

		handler := NewHandler(mockRepo)
		res, err := handler.UpdateGlobalCollege(ctx, collegeID, input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "global college not found")
	})
}
