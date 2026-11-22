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

// Unit tests for the UploadVideo handler logic without HTTP or AWS dependencies.
func TestHandler_UploadVideo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := models.PresignUploadRequestBody{
		OriginalFilename: "common-app-tips.mp4",
		ContentType:      "video/mp4",
		ContentLength:    1024,
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expected := &models.PresignUploadResponse{
			S3Key:     "videos/uuid-common-app-tips.mp4",
			UploadURL: "https://bucket.s3.amazonaws.com/videos/uuid-common-app-tips.mp4?X-Amz-...",
		}

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("PresignUpload", mock.Anything, input).Return(expected, nil)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.UploadVideo(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("PresignUpload", mock.Anything, input).
			Return(nil, errors.New("s3client: failed to presign upload URL"))

		handler := NewHandler(mockRepo, nil)
		res, err := handler.UploadVideo(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "s3client: failed to presign upload URL")
	})
}
