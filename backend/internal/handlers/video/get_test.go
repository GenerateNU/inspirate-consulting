package video

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the GetVideo handler, focused on which S3 keys each caller may fetch.
func TestHandler_GetVideo(t *testing.T) {
	t.Parallel()

	s3Key := "videos/uuid-common-app-tips.mp4"
	studentID := uuid.New()
	studentCtx := auth.WithStudentID(context.Background(), studentID)
	counselorCtx := auth.WithCounselorID(context.Background(), uuid.New())

	expected := &models.Video{
		S3Key:       s3Key,
		DownloadURL: "https://bucket.s3.amazonaws.com/videos/uuid-common-app-tips.mp4?X-Amz-...",
	}

	t.Run("counselor can fetch any key without an access check", func(t *testing.T) {
		t.Parallel()

		// The access mock has no expectations, so any access lookup fails the test
		mockRepo := mocks.NewVideoRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GetVideo", mock.Anything, s3Key).Return(expected, nil)

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetVideo(counselorCtx, s3Key)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("student can fetch a key that belongs to media they have access to", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockAccess.On("HasS3KeyAccess", mock.Anything, studentID.String(), s3Key).Return(true, nil)
		mockRepo.On("GetVideo", mock.Anything, s3Key).Return(expected, nil)

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetVideo(studentCtx, s3Key)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("student without access gets 404 and no download url is generated", func(t *testing.T) {
		t.Parallel()

		// No GetVideo expectation: presigning a URL would fail the test
		mockRepo := mocks.NewVideoRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockAccess.On("HasS3KeyAccess", mock.Anything, studentID.String(), s3Key).Return(false, nil)

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetVideo(studentCtx, s3Key)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	})

	t.Run("caller with no student or counselor profile is forbidden", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetVideo(context.Background(), s3Key)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("access check error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockAccess.On("HasS3KeyAccess", mock.Anything, studentID.String(), s3Key).Return(false, errors.New("database error"))

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetVideo(studentCtx, s3Key)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GetVideo", mock.Anything, s3Key).
			Return(nil, errors.New("s3client: failed to presign download URL"))

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetVideo(counselorCtx, s3Key)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "s3client: failed to presign download URL")
	})
}
