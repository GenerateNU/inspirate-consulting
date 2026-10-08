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
)

// Unit tests for the CreateEssayGroup handler logic without HTTP or database dependencies.
func TestHandler_CreateEssayGroup(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	studentID := uuid.New()
	description := "Essays for the Common App"

	input := &models.CreateEssayGroupInput{
		Body: models.CreateEssayGroupBody{
			Name:        "Common App",
			Description: &description,
			StudentID:   studentID,
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		createdGroup := &models.EssayGroups{
			ID:          uuid.New(),
			Name:        "Common App",
			Description: &description,
			StudentID:   studentID,
		}

		mockRepo := mocks.NewEssayGroupRepository(t)
		mockRepo.On(
			"CreateEssayGroup",
			mock.Anything,
			input.Body,
		).Return(createdGroup, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.CreateEssayGroup(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, &models.CreateEssayGroupOutput{
			Body: *createdGroup,
		}, res)
	})

	// Description is nullable in the database, so a nil pointer has to survive the handler.
	t.Run("success with no description", func(t *testing.T) {
		t.Parallel()

		noDescriptionInput := &models.CreateEssayGroupInput{
			Body: models.CreateEssayGroupBody{
				Name:      "Supplementals",
				StudentID: studentID,
			},
		}

		createdGroup := &models.EssayGroups{
			ID:        uuid.New(),
			Name:      "Supplementals",
			StudentID: studentID,
		}

		mockRepo := mocks.NewEssayGroupRepository(t)
		mockRepo.On(
			"CreateEssayGroup",
			mock.Anything,
			noDescriptionInput.Body,
		).Return(createdGroup, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.CreateEssayGroup(ctx, noDescriptionInput)

		assert.NoError(t, err)
		assert.Nil(t, res.Body.Description)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayGroupRepository(t)
		mockRepo.On(
			"CreateEssayGroup",
			mock.Anything,
			input.Body,
		).Return((*models.EssayGroups)(nil), errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.CreateEssayGroup(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
