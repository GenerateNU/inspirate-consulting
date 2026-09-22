package todoitem

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the GetTodoItemsByStudent handler logic without HTTP or database dependencies.
func TestHandler_GetTodoItemsByStudent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	studentID := uuid.NewString()
	userID := uuid.NewString()

	expectedOutput := []models.TodoItem{
		{
			ID:              uuid.NewString(),
			StudentID:       studentID,
			UserID:          userID,
			TodoDescription: "Submit the FAFSA",
		},
		{
			ID:              uuid.NewString(),
			StudentID:       studentID,
			UserID:          userID,
			TodoDescription: "Request a recommendation letter",
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemsByStudent", mock.Anything, studentID).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.GetTodoItemsByStudent(ctx)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
		assert.Len(t, res, 2)
	})

	t.Run("student with no todo items", func(t *testing.T) {
		t.Parallel()

		// A student with no tasks is an empty list, not an error
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemsByStudent", mock.Anything, studentID).Return([]models.TodoItem{}, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.GetTodoItemsByStudent(ctx)

		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemsByStudent", mock.Anything, studentID).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.GetTodoItemsByStudent(ctx)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
