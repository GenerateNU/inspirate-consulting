package video

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_DeleteVideo(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s3Key := "videos/uuid-common-app-tips.mp4"

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("DeleteVideo", mock.Anything, s3Key).Return(nil)

		handler := NewHandler(mockRepo)
		err := handler.DeleteVideo(ctx, s3Key)

		assert.NoError(t, err)
	})

	t.Run("not found error propagates", func(t *testing.T) {
		t.Parallel()

		notFound := errs.NotFound("video", "s3Key", s3Key)
		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("DeleteVideo", mock.Anything, s3Key).Return(notFound)

		handler := NewHandler(mockRepo)
		err := handler.DeleteVideo(ctx, s3Key)

		assert.Error(t, err)
		assert.Equal(t, notFound, err)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("DeleteVideo", mock.Anything, s3Key).
			Return(errors.New("s3client: failed to delete object"))

		handler := NewHandler(mockRepo)
		err := handler.DeleteVideo(ctx, s3Key)

		assert.Error(t, err)
		assert.EqualError(t, err, "s3client: failed to delete object")
	})
}
