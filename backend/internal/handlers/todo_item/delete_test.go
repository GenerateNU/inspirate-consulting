package todoitem

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

// Unit tests for the DeleteTodoItem handler logic without HTTP or database dependencies.
func TestHandler_DeleteTodoItem(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	id := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// The handler should pass the id straight through to the repository
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("DeleteTodoItem", mock.Anything, id).Return(nil)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		err := handler.DeleteTodoItem(ctx, id)

		// Verify result
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		// Simulate deleting an item that does not exist
		notFound := errs.NotFound("todo item", "id", id)
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("DeleteTodoItem", mock.Anything, id).Return(notFound)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		err := handler.DeleteTodoItem(ctx, id)

		// Verify error propagation
		assert.Error(t, err)
		assert.Equal(t, notFound, err)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("DeleteTodoItem", mock.Anything, id).Return(errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		err := handler.DeleteTodoItem(ctx, id)

		// Verify error propagation
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}
