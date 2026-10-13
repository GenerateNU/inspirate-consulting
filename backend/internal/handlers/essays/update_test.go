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

		updatedEssay := &models.Essays{
			ID:     essayID,
			Status: models.Submitted,
		}

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On(
			"UpdateStatus",
			mock.Anything,
			essayID,
			models.Submitted,
		).Return(updatedEssay, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.UpdateStatus(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, &models.UpdateStatusOutput{
			Body: models.EssayBody{Essay: updatedEssay},
		}, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On(
			"UpdateStatus",
			mock.Anything,
			essayID,
			models.Submitted,
		).Return((*models.Essays)(nil), errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.UpdateStatus(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
