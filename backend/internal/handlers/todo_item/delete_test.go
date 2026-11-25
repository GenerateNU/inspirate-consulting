package todoitem

import (
	"context"
	"errors"
	"net/http"
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
	ownerID := uuid.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Nil Authorizer allows everything, so this covers the happy path only
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemStudentID", mock.Anything, id).Return(ownerID, nil)
		mockRepo.On("DeleteTodoItem", mock.Anything, id).Return(nil)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		err := handler.DeleteTodoItem(ctx, id)

		// Verify result
		assert.NoError(t, err)
	})

	t.Run("student can delete their own todo", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemStudentID", mock.Anything, id).Return(ownerID, nil)
		mockRepo.On("DeleteTodoItem", mock.Anything, id).Return(nil)

		handler := NewHandler(mockRepo, realAuthorizer())
		err := handler.DeleteTodoItem(studentRoleCtx(ownerID), id)

		assert.NoError(t, err)
	})

	t.Run("student cannot delete another student's todo", func(t *testing.T) {
		t.Parallel()

		// No DeleteTodoItem expectation: the todo (and its media access) must not be deleted
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemStudentID", mock.Anything, id).Return(ownerID, nil)

		handler := NewHandler(mockRepo, realAuthorizer())
		err := handler.DeleteTodoItem(studentRoleCtx(uuid.New()), id)

		var httpErr errs.HTTPError
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		// A missing todo is reported by the owner lookup, before any delete
		notFound := errs.NotFound("todo items", "id", id)
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemStudentID", mock.Anything, id).Return(uuid.Nil, notFound)

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
		mockRepo.On("GetTodoItemStudentID", mock.Anything, id).Return(ownerID, nil)
		mockRepo.On("DeleteTodoItem", mock.Anything, id).Return(errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		err := handler.DeleteTodoItem(ctx, id)

		// Verify error propagation
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}
