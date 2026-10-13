package essay_groups

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

// Unit tests for the GetEssayGroups handler logic without HTTP or database dependencies.
func TestHandler_GetEssayGroups(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	studentID := uuid.New()

	input := &models.GetEssayGroupsInput{StudentID: studentID}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		groups := []models.EssayGroups{
			{ID: uuid.New(), Name: "Common App", StudentID: studentID},
			{ID: uuid.New(), Name: "Supplementals", StudentID: studentID},
		}

		mockRepo := mocks.NewEssayGroupRepository(t)
		mockRepo.On("ListEssayGroups", mock.Anything, studentID).Return(groups, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.GetEssayGroups(ctx, input)

		assert.NoError(t, err)
		require.Len(t, res.Body.EssayGroups, 2)
		assert.Equal(t, groups, res.Body.EssayGroups)
	})

	t.Run("student with no groups", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayGroupRepository(t)
		mockRepo.On("ListEssayGroups", mock.Anything, studentID).Return([]models.EssayGroups{}, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.GetEssayGroups(ctx, input)

		assert.NoError(t, err)
		assert.Empty(t, res.Body.EssayGroups)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayGroupRepository(t)
		mockRepo.On("ListEssayGroups", mock.Anything, studentID).
			Return([]models.EssayGroups(nil), errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.GetEssayGroups(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
