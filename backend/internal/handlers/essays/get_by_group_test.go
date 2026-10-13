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

// Unit tests for the GetEssaysByGroup handler logic without HTTP or database dependencies.
func TestHandler_GetEssaysByGroup(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	groupID := uuid.New()

	input := &models.GetEssaysByGroupInput{EssayGroupID: groupID}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		essays := []models.Essays{
			{ID: uuid.New(), Type: "personal-statement", EssayGroupID: &groupID},
			{ID: uuid.New(), Type: "supplemental", EssayGroupID: &groupID},
		}

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssaysByGroup", mock.Anything, groupID).Return(essays, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.GetEssaysByGroup(ctx, input)

		assert.NoError(t, err)
		require.Len(t, res.Body.Essays, 2)
		assert.Equal(t, essays, res.Body.Essays)
	})

	t.Run("group with no essays", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssaysByGroup", mock.Anything, groupID).Return([]models.Essays{}, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.GetEssaysByGroup(ctx, input)

		assert.NoError(t, err)
		assert.Empty(t, res.Body.Essays)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssaysByGroup", mock.Anything, groupID).
			Return([]models.Essays(nil), errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.GetEssaysByGroup(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
