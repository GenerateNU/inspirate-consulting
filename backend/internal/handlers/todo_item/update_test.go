package todoitem

import (
	"context"
	"errors"
	"testing"
	"time"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the UpdateTodoItemCompletedAt handler logic without HTTP or database dependencies.
// The handler translates the frontend's completed flag into a completed_at timestamp.
func TestHandler_UpdateTodoItemCompletedAt(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()

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
		mockRepo.On("UpdateTodoItemCompletedAt", mock.Anything, id, mock.MatchedBy(func(ts *time.Time) bool {
			return ts != nil && time.Since(*ts) < time.Minute
		})).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo)
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
		mockRepo.On("UpdateTodoItemCompletedAt", mock.Anything, id, (*time.Time)(nil)).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo)
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
		mockRepo.On("UpdateTodoItemCompletedAt", mock.Anything, id, mock.Anything).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.UpdateTodoItemCompletedAt(ctx, id, true)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
