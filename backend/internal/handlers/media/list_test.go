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

// Unit tests for the ListAllMedia handler logic without HTTP or database dependencies.
func TestHandler_ListAllMedia(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	expectedOutput := []models.Media{
		{
			ID:           uuid.NewString(),
			Title:        "Common App Essay Tips",
			LengthInMins: 15,
			S3Key:        "video-bucket/common-app-tips.mp4",
		},
		{
			ID:           uuid.NewString(),
			Title:        "Navigating the FAFSA",
			LengthInMins: 20,
			S3Key:        "video-bucket/fafsa.mp4",
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("ListAllMedia", mock.Anything).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.ListAllMedia(ctx)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
		assert.Len(t, res, 2)
	})

	t.Run("no media", func(t *testing.T) {
		t.Parallel()

		// No media uploaded yet is an empty list, not an error
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("ListAllMedia", mock.Anything).Return([]models.Media{}, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.ListAllMedia(ctx)

		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("ListAllMedia", mock.Anything).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.ListAllMedia(ctx)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
