package user

import (
	"context"
	"errors"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the CreateGreeting handler logic without HTTP or database dependencies.
func TestHandler_CreateGreeting(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := &models.CreateUserInput{
		Body: models.GreetingRequestBody{
			Name: "Alice",
		},
	}
	expectedOutput := &models.CreateGreetingOutput{
		Body: models.GreetingMessageBody{
			Message: "Hello, Alice!",
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewGreetingRepository(t)
		mockRepo.On("CreateGreeting", mock.Anything, *input).Return(expectedOutput, nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.CreateUser(ctx, input)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewGreetingRepository(t)
		mockRepo.On("CreateGreeting", mock.Anything, *input).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.CreateGreeting(ctx, input)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
