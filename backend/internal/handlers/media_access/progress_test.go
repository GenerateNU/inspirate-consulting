package mediaaccess

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the RecordMediaProgress handler logic without HTTP or database dependencies.
func TestHandler_RecordMediaProgress(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	mediaID := uuid.NewString()
	ctx := auth.WithStudentID(context.Background(), studentID)
	now := time.Now()

	t.Run("started records started only", func(t *testing.T) {
		t.Parallel()

		expected := &models.MediaAccess{ID: uuid.NewString(), StudentID: studentID.String(), MediaID: mediaID, StartedAt: &now}
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RecordMediaProgress", mock.Anything, studentID.String(), mediaID, false).Return(expected, nil)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.RecordMediaProgress(ctx, mediaID, models.MediaProgressStarted)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("completed records completion", func(t *testing.T) {
		t.Parallel()

		expected := &models.MediaAccess{ID: uuid.NewString(), StudentID: studentID.String(), MediaID: mediaID, StartedAt: &now, CompletedAt: &now}
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RecordMediaProgress", mock.Anything, studentID.String(), mediaID, true).Return(expected, nil)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.RecordMediaProgress(ctx, mediaID, models.MediaProgressCompleted)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("progress is always recorded for the calling student", func(t *testing.T) {
		t.Parallel()

		// The student id comes from the auth context, never from the request,
		// so a student can't record progress on someone else's access row
		otherStudentID := uuid.New()
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RecordMediaProgress", mock.Anything, otherStudentID.String(), mediaID, false).
			Return(&models.MediaAccess{StudentID: otherStudentID.String()}, nil)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.RecordMediaProgress(auth.WithStudentID(context.Background(), otherStudentID), mediaID, models.MediaProgressStarted)

		assert.NoError(t, err)
		assert.Equal(t, otherStudentID.String(), res.StudentID)
	})

	t.Run("caller without a student profile is forbidden", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.RecordMediaProgress(auth.WithCounselorID(context.Background(), uuid.New()), mediaID, models.MediaProgressStarted)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("unknown status is a bad request", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.RecordMediaProgress(ctx, mediaID, models.MediaProgressStatus("paused"))

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("no access row is a 404 from the repository", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RecordMediaProgress", mock.Anything, studentID.String(), mediaID, false).
			Return(nil, errs.NotFound("media", "id", mediaID))

		handler := NewHandler(mockRepo, nil)
		res, err := handler.RecordMediaProgress(ctx, mediaID, models.MediaProgressStarted)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RecordMediaProgress", mock.Anything, studentID.String(), mediaID, true).
			Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo, nil)
		res, err := handler.RecordMediaProgress(ctx, mediaID, models.MediaProgressCompleted)

		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
