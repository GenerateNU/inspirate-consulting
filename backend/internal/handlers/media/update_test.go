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

// Unit tests for the UpdateMedia handler logic without HTTP or database dependencies.
func TestHandler_UpdateMedia(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()
	newTitle := "Common App Essay Tips (2026)"

	expectedOutput := &models.Media{
		ID:           id,
		Title:        newTitle,
		Description:  "How to write a strong personal statement",
		LengthInMins: 15,
		S3Key:        "video-bucket/common-app-tips.mp4",
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("UpdateMedia", mock.Anything, id, newTitle).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.UpdateMedia(ctx, id, newTitle)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("UpdateMedia", mock.Anything, id, newTitle).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.UpdateMedia(ctx, id, newTitle)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
