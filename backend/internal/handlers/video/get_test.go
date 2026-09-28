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

func TestHandler_GetVideo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s3Key := "videos/uuid-common-app-tips.mp4"

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expected := &models.Video{
			S3Key: s3Key,
			DownloadURL: "https://bucket.s3.amazonaws.com/videos/uuid-common-app-tips.mp4?X-Amz-...",
		}

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("GetVideo", mock.Anything, s3Key).Return(expected, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.GetVideo(ctx, s3Key)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("GetVideo", mock.Anything, s3Key).
			Return(nil, errors.New("s3client: failed to presign download URL"))

		handler := NewHandler(mockRepo)
		res, err := handler.GetVideo(ctx, s3Key)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "s3client: failed to presign download URL")
	})
}