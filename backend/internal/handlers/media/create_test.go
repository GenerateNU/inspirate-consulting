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

// Unit tests for the CreateMedia handler logic without HTTP or database dependencies.
func TestHandler_CreateMedia(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schoolYear := 11

	input := &models.CreateMediaRequestBody{
		Title:        "Common App Essay Tips",
		Description:  "How to write a strong personal statement",
		LengthInMins: 15,
		SchoolYear:   &schoolYear,
		S3Key:        "video-bucket/common-app-tips.mp4",
	}
	expectedOutput := &models.Media{
		ID:           uuid.NewString(),
		Title:        "Common App Essay Tips",
		Description:  "How to write a strong personal statement",
		LengthInMins: 15,
		SchoolYear:   &schoolYear,
		S3Key:        "video-bucket/common-app-tips.mp4",
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("CreateMedia", mock.Anything, input).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.CreateMedia(ctx, input)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("media without school year", func(t *testing.T) {
		t.Parallel()

		// Media not tied to a specific school year should come back with a nil school year
		noYearInput := &models.CreateMediaRequestBody{
			Title:        "Intro to College Applications",
			LengthInMins: 10,
			S3Key:        "video-bucket/intro.mp4",
		}
		noYearOutput := &models.Media{
			ID:           uuid.NewString(),
			Title:        "Intro to College Applications",
			LengthInMins: 10,
			S3Key:        "video-bucket/intro.mp4",
		}

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("CreateMedia", mock.Anything, noYearInput).Return(noYearOutput, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.CreateMedia(ctx, noYearInput)

		assert.NoError(t, err)
		assert.Nil(t, res.SchoolYear)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("CreateMedia", mock.Anything, input).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.CreateMedia(ctx, input)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
