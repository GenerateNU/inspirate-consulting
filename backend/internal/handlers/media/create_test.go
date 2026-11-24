package media

import (
	"context"
	"errors"
	"net/http"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the CreateMedia handler logic without HTTP or database dependencies.
func TestHandler_CreateMedia(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	noYoutubeID := (*string)(nil)

	schoolYear := models.Junior
	input := &models.CreateMediaRequestBody{
		Title:        "Common App Essay Tips",
		Description:  "How to write a strong personal statement",
		LengthInMins: ptr(15),
		SchoolYear:   &schoolYear,
		S3Key:        ptr("video-bucket/common-app-tips.mp4"),
	}
	expectedOutput := &models.Media{
		ID:           uuid.NewString(),
		Title:        "Common App Essay Tips",
		Description:  "How to write a strong personal statement",
		LengthInMins: ptr(15),
		SchoolYear:   &schoolYear,
		S3Key:        ptr("video-bucket/common-app-tips.mp4"),
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("CreateMedia", mock.Anything, input, noYoutubeID).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
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
			LengthInMins: ptr(10),
			S3Key:        ptr("video-bucket/intro.mp4"),
		}
		noYearOutput := &models.Media{
			ID:           uuid.NewString(),
			Title:        "Intro to College Applications",
			LengthInMins: ptr(10),
			S3Key:        ptr("video-bucket/intro.mp4"),
		}

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("CreateMedia", mock.Anything, noYearInput, noYoutubeID).Return(noYearOutput, nil)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.CreateMedia(ctx, noYearInput)

		assert.NoError(t, err)
		assert.Nil(t, res.SchoolYear)
	})

	t.Run("youtube url is parsed to a video id", func(t *testing.T) {
		t.Parallel()

		youtubeInput := &models.CreateMediaRequestBody{
			Title:      "FAFSA Walkthrough",
			YoutubeURL: ptr("https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=42s"),
		}
		youtubeOutput := &models.Media{
			ID:             uuid.NewString(),
			Title:          "FAFSA Walkthrough",
			YoutubeVideoID: ptr("dQw4w9WgXcQ"),
		}

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("CreateMedia", mock.Anything, youtubeInput, ptr("dQw4w9WgXcQ")).Return(youtubeOutput, nil)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.CreateMedia(ctx, youtubeInput)

		assert.NoError(t, err)
		assert.Equal(t, youtubeOutput, res)
	})

	badRequests := map[string]*models.CreateMediaRequestBody{
		"neither s3_key nor youtube_url": {Title: "No source"},
		"both s3_key and youtube_url": {
			Title:      "Two sources",
			S3Key:      ptr("video-bucket/intro.mp4"),
			YoutubeURL: ptr("https://youtu.be/dQw4w9WgXcQ"),
		},
		"non-youtube url": {
			Title:      "Not YouTube",
			YoutubeURL: ptr("https://evil.example.com/watch?v=dQw4w9WgXcQ"),
		},
	}
	for name, badInput := range badRequests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Repository must not be called; the mock fails the test on any unexpected call
			mockRepo := mocks.NewMediaRepository(t)

			handler := NewHandler(mockRepo, nil)
			res, err := handler.CreateMedia(ctx, badInput)

			assert.Nil(t, res)
			var httpErr errs.HTTPError
			assert.ErrorAs(t, err, &httpErr)
			assert.Equal(t, http.StatusBadRequest, httpErr.Code)
		})
	}

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("CreateMedia", mock.Anything, input, noYoutubeID).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		res, err := handler.CreateMedia(ctx, input)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
