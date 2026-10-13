package media

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the GetMedia handler logic without HTTP or database dependencies.
func TestHandler_GetMedia(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()

	expectedOutput := &models.Media{
		ID:           id,
		Title:        "Common App Essay Tips",
		Description:  "How to write a strong personal statement",
		LengthInMins: 15,
		S3Key:        "video-bucket/common-app-tips.mp4",
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("GetMedia", mock.Anything, id).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.GetMedia(ctx, id)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("GetMedia", mock.Anything, id).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.GetMedia(ctx, id)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
