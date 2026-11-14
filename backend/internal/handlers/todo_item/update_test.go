package todoitem

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Unit tests for the UpdateTodoItemCompletedAt handler logic without HTTP or database dependencies.
// The handler translates the frontend's completed flag into a completed_at timestamp.
func TestHandler_UpdateTodoItemCompletedAt(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()
	studentID := uuid.New()

	t.Run("completed true sets completed_at", func(t *testing.T) {
		t.Parallel()

		completedAt := time.Now()
		expectedOutput := &models.TodoItem{
			ID:              id,
			TodoDescription: "Submit the FAFSA",
			CompletedAt:     &completedAt,
		}

		// completed=true should hand the repository a non-nil, current timestamp
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemStudentID", mock.Anything, id).Return(studentID, nil)
		mockRepo.On("UpdateTodoItemCompletedAt", mock.Anything, id, mock.MatchedBy(func(ts *time.Time) bool {
			return ts != nil && time.Since(*ts) < time.Minute
		})).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		res, err := handler.UpdateTodoItemCompletedAt(ctx, id, true)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
		assert.NotNil(t, res.CompletedAt)
	})

	t.Run("completed false clears completed_at", func(t *testing.T) {
		t.Parallel()

		expectedOutput := &models.TodoItem{
			ID:              id,
			TodoDescription: "Completed by mistake",
			CompletedAt:     nil,
		}

		// completed=false should hand the repository a nil timestamp, which reopens the item
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemStudentID", mock.Anything, id).Return(studentID, nil)
		mockRepo.On("UpdateTodoItemCompletedAt", mock.Anything, id, (*time.Time)(nil)).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		res, err := handler.UpdateTodoItemCompletedAt(ctx, id, false)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
		assert.Nil(t, res.CompletedAt)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewTodoItemRepository(t)
		// the ownership lookup must succeed so the failure comes from the update step
		mockRepo.On("GetTodoItemStudentID", mock.Anything, id).Return(studentID, nil)
		mockRepo.On("UpdateTodoItemCompletedAt", mock.Anything, id, mock.Anything).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		res, err := handler.UpdateTodoItemCompletedAt(ctx, id, true)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})

	t.Run("todo item not found", func(t *testing.T) {
		t.Parallel()

		// a failed ownership lookup must stop the handler before any update
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemStudentID", mock.Anything, id).
			Return(uuid.Nil, errs.NotFound("todo items", "id", id))

		handler := NewHandler(mockRepo, nil)
		res, err := handler.UpdateTodoItemCompletedAt(ctx, id, true)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusNotFound, httpErr.GetStatus())
		mockRepo.AssertNotCalled(t, "UpdateTodoItemCompletedAt", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("student cannot update another student's todo item", func(t *testing.T) {
		t.Parallel()

		// the item belongs to a different student, so the real authorizer must reject it
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemStudentID", mock.Anything, id).Return(uuid.New(), nil)

		handler := NewHandler(mockRepo, realAuthorizer())
		res, err := handler.UpdateTodoItemCompletedAt(studentRoleCtx(studentID), id, true)

		assert.Nil(t, res)
		var httpErr errs.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.GetStatus())
		mockRepo.AssertNotCalled(t, "UpdateTodoItemCompletedAt", mock.Anything, mock.Anything, mock.Anything)
	})
}

// Unit tests for the UpdateTodoItemDescription handler logic without HTTP or database dependencies.
func TestHandler_UpdateTodoItemDescription(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expectedOutput := &models.TodoItem{
			ID:              id,
			TodoDescription: "Submit the CSS Profile",
		}

		// the new description should be passed through to the repository unchanged
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("UpdateTodoItemDescription", mock.Anything, id, "Submit the CSS Profile").Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		res, err := handler.UpdateTodoItemDescription(ctx, id, "Submit the CSS Profile")

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("UpdateTodoItemDescription", mock.Anything, id, mock.Anything).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		res, err := handler.UpdateTodoItemDescription(ctx, id, "Submit the CSS Profile")

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}

// Unit tests for the UpdateTodoItemDeadline handler logic without HTTP or database dependencies.
// The handler rejects deadlines in the past and passes everything else to the repository.
func TestHandler_UpdateTodoItemDeadline(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()

	t.Run("future deadline is saved", func(t *testing.T) {
		t.Parallel()

		deadline := time.Now().Add(7 * 24 * time.Hour)
		expectedOutput := &models.TodoItem{
			ID:              id,
			TodoDescription: "Submit the FAFSA",
			Deadline:        &deadline,
		}

		// a future deadline should be passed through to the repository unchanged
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("UpdateTodoItemDeadline", mock.Anything, id, &deadline).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		res, err := handler.UpdateTodoItemDeadline(ctx, id, &deadline)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("nil deadline clears it", func(t *testing.T) {
		t.Parallel()

		expectedOutput := &models.TodoItem{
			ID:              id,
			TodoDescription: "Submit the FAFSA",
			Deadline:        nil,
		}

		// a nil deadline should reach the repository as nil, which removes the deadline
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("UpdateTodoItemDeadline", mock.Anything, id, (*time.Time)(nil)).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		res, err := handler.UpdateTodoItemDeadline(ctx, id, nil)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
		assert.Nil(t, res.Deadline)
	})

	t.Run("past deadline is rejected", func(t *testing.T) {
		t.Parallel()

		// Mock should NOT be called when the deadline is in the past
		mockRepo := mocks.NewTodoItemRepository(t)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		past := time.Now().Add(-24 * time.Hour)
		res, err := handler.UpdateTodoItemDeadline(ctx, id, &past)

		// Verify the request is rejected before reaching the repository
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "deadline cannot be in the past")
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("UpdateTodoItemDeadline", mock.Anything, id, mock.Anything).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		deadline := time.Now().Add(24 * time.Hour)
		res, err := handler.UpdateTodoItemDeadline(ctx, id, &deadline)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
