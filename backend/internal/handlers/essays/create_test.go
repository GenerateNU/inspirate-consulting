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

// Unit tests for the CreateEssay handler logic without HTTP or database dependencies.
func TestHandler_CreateEssay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	studentID := uuid.New()
	collegeID := int64(42)

	input := &models.CreateEssayInput{
		Body: models.CreateEssayBody{
			StudentID:     studentID,
			Type:          "personal-statement",
			CollegeID:     &collegeID,
			LinkToContent: "https://docs.google.com/document/d/abc123",
		},
	}
	// The row struct the handler is expected to build from that input. ID and
	// Status stay zero so the database applies its own defaults.
	expectedEssay := models.Essays{
		StudentID:     studentID,
		Type:          "personal-statement",
		CollegeID:     &collegeID,
		LinkToContent: "https://docs.google.com/document/d/abc123",
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Set up mock repository expectation
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("CreateEssay", mock.Anything, expectedEssay).Return(nil)

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.CreateEssay(ctx, input)

		// Verify result
		assert.NoError(t, err)
		assert.Equal(t, &models.CreateEssayOutput{}, res)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate database repository failure
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("CreateEssay", mock.Anything, expectedEssay).Return(errors.New("database error"))

		// Execute handler
		handler := NewHandler(mockRepo)
		res, err := handler.CreateEssay(ctx, input)

		// Verify error propagation
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
