package essays

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

// Unit tests for the UpdateStatus handler logic without HTTP or database dependencies.
func TestHandler_UpdateStatus(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	essayID := uuid.New()

	input := &models.UpdateStatusInput{
		EssayID: essayID,
		Body: models.UpdateStatusBody{
			Status: models.Submitted,
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation: the ID comes from the path and
		// the status from the body, so both must reach the repository
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("UpdateStatus", mock.Anything, essayID, models.Submitted).Return(nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.UpdateStatus(ctx, input)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, &models.UpdateStatusOutput{}, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("UpdateStatus", mock.Anything, essayID, models.Submitted).Return(errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.UpdateStatus(ctx, input)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
