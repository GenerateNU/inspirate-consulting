package todoitem

import (
	"context"
	"errors"
	"testing"
	"time"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the CreateTodoItem handler logic without HTTP or database dependencies.
func TestHandler_CreateTodoItem(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	userID := uuid.New()
	ctx := auth.WithUserID(auth.WithStudentID(context.Background(), studentID), userID)
	deadline := time.Now().Add(7 * 24 * time.Hour)

	input := &models.CreateTodoItemRequestBody{
		TodoDescription: "Finish the Common App essay",
		Deadline:        &deadline,
	}
	expectedOutput := &models.TodoItem{
		ID:              uuid.NewString(),
		StudentID:       studentID.String(),
		UserID:          userID.String(),
		TodoDescription: "Finish the Common App essay",
		Deadline:        &deadline,
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("CreateTodoItem", mock.Anything, studentID.String(), userID.String(), *input).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		res, err := handler.CreateTodoItem(ctx, input)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("new item is not completed", func(t *testing.T) {
		t.Parallel()

		// A freshly created todo item should come back with no completion timestamp
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("CreateTodoItem", mock.Anything, studentID.String(), userID.String(), *input).Return(expectedOutput, nil)

		handler := NewHandler(mockRepo, nil)
		res, err := handler.CreateTodoItem(ctx, input)

		assert.NoError(t, err)
		assert.Nil(t, res.CompletedAt)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("CreateTodoItem", mock.Anything, studentID.String(), userID.String(), *input).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo, nil)
		res, err := handler.CreateTodoItem(ctx, input)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
