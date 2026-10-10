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
	"github.com/stretchr/testify/require"
)

// Unit tests for the UpdateEssayGroup handler logic without HTTP or database dependencies.
func TestHandler_UpdateEssayGroup(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	essayID := uuid.New()
	groupID := uuid.New()

	t.Run("moves the essay into a group", func(t *testing.T) {
		t.Parallel()

		updatedEssay := &models.Essays{
			ID:           essayID,
			EssayGroupID: &groupID,
		}

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("UpdateEssayGroup", mock.Anything, essayID, &groupID).Return(updatedEssay, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.UpdateEssayGroup(ctx, &models.UpdateEssayGroupInput{
			EssayID: essayID,
			Body:    models.UpdateEssayGroupBody{EssayGroupID: &groupID},
		})

		assert.NoError(t, err)
		assert.Equal(t, &models.UpdateEssayGroupOutput{
			Body: models.EssayBody{Essay: updatedEssay},
		}, res)
	})

	// A null essay_group_id is how a client removes an essay from its group.
	t.Run("removes the essay from its group", func(t *testing.T) {
		t.Parallel()

		updatedEssay := &models.Essays{
			ID:           essayID,
			EssayGroupID: nil,
		}

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("UpdateEssayGroup", mock.Anything, essayID, (*uuid.UUID)(nil)).Return(updatedEssay, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.UpdateEssayGroup(ctx, &models.UpdateEssayGroupInput{
			EssayID: essayID,
			Body:    models.UpdateEssayGroupBody{EssayGroupID: nil},
		})

		assert.NoError(t, err)
		require.NotNil(t, res.Body.Essay)
		assert.Nil(t, res.Body.Essay.EssayGroupID)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("UpdateEssayGroup", mock.Anything, essayID, &groupID).
			Return((*models.Essays)(nil), errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.UpdateEssayGroup(ctx, &models.UpdateEssayGroupInput{
			EssayID: essayID,
			Body:    models.UpdateEssayGroupBody{EssayGroupID: &groupID},
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
