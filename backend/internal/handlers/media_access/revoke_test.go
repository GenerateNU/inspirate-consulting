package mediaaccess

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the RevokeMediaAccess handler logic without HTTP or database dependencies.
func TestHandler_RevokeMediaAccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()
	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectations: look up the owning student, then revoke
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GetMediaAccessStudentID", mock.Anything, id).Return(studentID, nil)
		mockRepo.On("RevokeMediaAccess", mock.Anything, id).Return(nil)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		err := handler.RevokeMediaAccess(ctx, id)

		// Verify result
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		// RevokeMediaAccess should NOT be called when the media access row does not exist
		notFound := errs.NotFound("media_access", "id", id)
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GetMediaAccessStudentID", mock.Anything, id).Return(uuid.Nil, notFound)

		handler := NewHandler(mockRepo, nil)
		err := handler.RevokeMediaAccess(ctx, id)

		// Verify error propagation
		assert.Error(t, err)
		assert.Equal(t, notFound, err)
		mockRepo.AssertNotCalled(t, "RevokeMediaAccess", mock.Anything, mock.Anything)
	})

	t.Run("student lookup error", func(t *testing.T) {
		t.Parallel()

		// Simulate a database failure while resolving the owning student
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GetMediaAccessStudentID", mock.Anything, id).Return(uuid.Nil, errors.New("database error"))

		handler := NewHandler(mockRepo, nil)
		err := handler.RevokeMediaAccess(ctx, id)

		assert.EqualError(t, err, "database error")
		mockRepo.AssertNotCalled(t, "RevokeMediaAccess", mock.Anything, mock.Anything)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure on revoke
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GetMediaAccessStudentID", mock.Anything, id).Return(studentID, nil)
		mockRepo.On("RevokeMediaAccess", mock.Anything, id).Return(errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		err := handler.RevokeMediaAccess(ctx, id)

		// Verify error propagation
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}
