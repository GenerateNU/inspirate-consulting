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

// Unit tests for the UpdateMedia handler logic without HTTP or database dependencies.
func TestHandler_UpdateMedia(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()

	expectedOutput := &models.Media{
		ID:           id,
		Title:        "Common App Essay Tips (2026)",
		Description:  "Updated description",
		LengthInMins: ptr(15),
		S3Key:        ptr("video-bucket/common-app-tips.mp4"),
	}

	t.Run("title and description pass through", func(t *testing.T) {
		t.Parallel()

		body := &models.UpdateMediaInputBody{
			Title:       ptr("Common App Essay Tips (2026)"),
			Description: ptr("Updated description"),
		}
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("UpdateMedia", mock.Anything, id, models.UpdateMediaFields{
			Title:       body.Title,
			Description: body.Description,
		}).Return(expectedOutput, nil)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.UpdateMedia(ctx, id, body)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("new s3_key passes through", func(t *testing.T) {
		t.Parallel()

		body := &models.UpdateMediaInputBody{S3Key: ptr("video-bucket/new.mp4")}
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("UpdateMedia", mock.Anything, id, models.UpdateMediaFields{S3Key: body.S3Key}).Return(expectedOutput, nil)

		handler := NewHandler(mockRepo, nil)
		_, err := handler.UpdateMedia(ctx, id, body)

		assert.NoError(t, err)
	})

	t.Run("youtube url is parsed to a video id", func(t *testing.T) {
		t.Parallel()

		youtubeOutput := &models.Media{ID: id, Title: "Common App Essay Tips", YoutubeVideoID: ptr("dQw4w9WgXcQ")}

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("UpdateMedia", mock.Anything, id, models.UpdateMediaFields{YoutubeVideoID: ptr("dQw4w9WgXcQ")}).Return(youtubeOutput, nil)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.UpdateMedia(ctx, id, &models.UpdateMediaInputBody{YoutubeURL: ptr("https://youtu.be/dQw4w9WgXcQ")})

		assert.NoError(t, err)
		assert.Equal(t, youtubeOutput, res)
	})

	badRequests := map[string]*models.UpdateMediaInputBody{
		"no fields":       {},
		"non-youtube url": {YoutubeURL: ptr("https://vimeo.com/123456")},
		"both s3_key and youtube_url": {
			S3Key:      ptr("video-bucket/new.mp4"),
			YoutubeURL: ptr("https://youtu.be/dQw4w9WgXcQ"),
		},
	}
	for name, body := range badRequests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Repository must not be called; the mock fails the test on any unexpected call
			mockRepo := mocks.NewMediaRepository(t)

			handler := NewHandler(mockRepo, nil)
			res, err := handler.UpdateMedia(ctx, id, body)

			assert.Nil(t, res)
			var httpErr errs.HTTPError
			assert.ErrorAs(t, err, &httpErr)
			assert.Equal(t, http.StatusBadRequest, httpErr.Code)
		})
	}

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		body := &models.UpdateMediaInputBody{Title: ptr("New Title")}
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("UpdateMedia", mock.Anything, id, models.UpdateMediaFields{Title: body.Title}).Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo, nil)
		res, err := handler.UpdateMedia(ctx, id, body)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
