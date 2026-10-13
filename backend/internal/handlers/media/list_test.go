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
	limit, offset := 20, 0

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
		mockRepo.On("ListAllMedia", mock.Anything, limit, offset).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.ListAllMedia(ctx, limit, offset)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
		assert.Len(t, res, 2)
	})

	t.Run("passes limit and offset through to repository", func(t *testing.T) {
		t.Parallel()

		// The mock only matches the exact pagination values, so a mismatch fails the test
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("ListAllMedia", mock.Anything, 1, 1).Return(expectedOutput[1:], nil)

		handler := NewHandler(mockRepo)
		res, err := handler.ListAllMedia(ctx, 1, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput[1:], res)
		assert.Len(t, res, 1)
	})

	t.Run("no media", func(t *testing.T) {
		t.Parallel()

		// No media uploaded yet is an empty list, not an error
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("ListAllMedia", mock.Anything, limit, offset).Return([]models.Media{}, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.ListAllMedia(ctx, limit, offset)

		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("ListAllMedia", mock.Anything, limit, offset).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.ListAllMedia(ctx, limit, offset)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
