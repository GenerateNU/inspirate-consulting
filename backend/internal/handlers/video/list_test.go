package video

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ListVideos(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expected := []models.Video{
			{S3Key: "videos/uuid-alpha.mp4", DownloadURL: "https://bucket.s3.amazonaws.com/videos/uuid-alpha.mp4?X-Amz-..."},
			{S3Key: "videos/uuid-beta.mp4", DownloadURL: "https://bucket.s3.amazonaws.com/videos/uuid-beta.mp4?X-Amz-..."},
		}

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("ListVideos", mock.Anything).Return(expected, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.ListVideos(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("ListVideos", mock.Anything).Return([]models.Video{}, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.ListVideos(ctx)

		assert.NoError(t, err)
		assert.Equal(t, []models.Video{}, res)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("ListVideos", mock.Anything).Return(nil, errors.New("s3client: failed to list objects"))

		handler := NewHandler(mockRepo)
		res, err := handler.ListVideos(ctx)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "s3client: failed to list objects")
	})
}
