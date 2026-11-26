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

// Unit tests for the GetMedia handler, focused on who is allowed to see which media.
func TestHandler_GetMedia(t *testing.T) {
	t.Parallel()

	id := uuid.NewString()
	studentID := uuid.New()

	expectedOutput := &models.Media{
		ID:           id,
		Title:        "Common App Essay Tips",
		Description:  "How to write a strong personal statement",
		LengthInMins: ptr(15),
		S3Key:        ptr("video-bucket/common-app-tips.mp4"),
	}

	t.Run("counselor can fetch any media without an access check", func(t *testing.T) {
		t.Parallel()

		// The access mock has no expectations, so any access lookup fails the test
		mockRepo := mocks.NewMediaRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GetMedia", mock.Anything, id).Return(expectedOutput, nil)

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetMedia(counselorCtx(uuid.New()), id)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("student with access can fetch the media", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockAccess.On("HasMediaAccess", mock.Anything, studentID.String(), id).Return(true, nil)
		mockRepo.On("GetMedia", mock.Anything, id).Return(expectedOutput, nil)

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetMedia(studentCtx(studentID), id)

		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("student without access gets 404 and the media is never fetched", func(t *testing.T) {
		t.Parallel()

		// No GetMedia expectation: fetching the media would fail the test
		mockRepo := mocks.NewMediaRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockAccess.On("HasMediaAccess", mock.Anything, studentID.String(), id).Return(false, nil)

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetMedia(studentCtx(studentID), id)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	})

	t.Run("caller with no student or counselor profile is forbidden", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetMedia(context.Background(), id)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("access check error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockAccess.On("HasMediaAccess", mock.Anything, studentID.String(), id).Return(false, errors.New("database error"))

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetMedia(studentCtx(studentID), id)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GetMedia", mock.Anything, id).Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo, mockAccess)
		res, err := handler.GetMedia(counselorCtx(uuid.New()), id)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
