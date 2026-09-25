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

// Unit tests for the GetEssaysFromStudent handler logic without HTTP or database dependencies.
func TestHandler_GetEssaysFromStudent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	studentID := uuid.New()
	collegeID := int64(42)

	input := &models.GetEssaysFromStudentInput{
		StudentID: studentID,
	}
	storedEssays := []models.Essays{
		{
			ID:            uuid.New(),
			StudentID:     studentID,
			Type:          "personal-statement",
			CollegeID:     &collegeID,
			LinkToContent: "https://docs.google.com/document/d/abc123",
			Status:        models.Draft,
		},
	}
	expectedOutput := &models.GetEssaysFromStudentOutput{
		Body: models.EssayListBody{
			Essays: storedEssays,
		},
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssaysFromStudent", mock.Anything, studentID).Return(storedEssays, nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.GetEssaysFromStudent(ctx, input)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssaysFromStudent", mock.Anything, studentID).Return(nil, errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.GetEssaysFromStudent(ctx, input)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
