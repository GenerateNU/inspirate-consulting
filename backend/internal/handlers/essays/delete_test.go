package essays

import (
	"context"
	"errors"
	"testing"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit tests for the DeleteEssay handler logic without HTTP or database dependencies.
func TestHandler_DeleteEssay(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	essayID := uuid.New()
	essay := &models.Essays{
		ID:            essayID,
		StudentID:     studentID,
		Type:          "personal-statement",
		LinkToContent: "https://docs.google.com/document/d/abc123",
		Status:        models.Draft,
	}

	t.Run("success - student deletes their essay", func(t *testing.T) {
		t.Parallel()

		ctx := auth.WithStudentID(context.Background(), studentID)

		// Set up mock repository expectations
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssay", mock.Anything, essayID).Return(essay, nil)
		mockRepo.On("DeleteEssay", mock.Anything, essayID, studentID).Return(nil)

		// Execute handler
		handler := NewHandler(mockRepo, nil, nil)
		err := handler.DeleteEssay(ctx, essayID)

		// Verify result
		assert.NoError(t, err)
	})

	t.Run("success - counselor deletes using the essay owner's student id", func(t *testing.T) {
		t.Parallel()

		// No student in context, so the owner must come from the essay itself
		ctx := auth.WithCounselorID(context.Background(), uuid.New())

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssay", mock.Anything, essayID).Return(essay, nil)
		mockRepo.On("DeleteEssay", mock.Anything, essayID, studentID).Return(nil)

		handler := NewHandler(mockRepo, nil, nil)
		err := handler.DeleteEssay(ctx, essayID)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		// DeleteEssay should NOT be called when the essay does not exist
		notFound := errs.NotFound("essay", "id", essayID.String())
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssay", mock.Anything, essayID).Return(nil, notFound)

		handler := NewHandler(mockRepo, nil, nil)
		err := handler.DeleteEssay(context.Background(), essayID)

		// Verify error propagation
		assert.Error(t, err)
		assert.Equal(t, notFound, err)
	})

	t.Run("repository error on delete", func(t *testing.T) {
		t.Parallel()

		ctx := auth.WithStudentID(context.Background(), studentID)

		// Simulate database repository failure
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssay", mock.Anything, essayID).Return(essay, nil)
		mockRepo.On("DeleteEssay", mock.Anything, essayID, studentID).Return(errors.New("database error"))

		handler := NewHandler(mockRepo, nil, nil)
		err := handler.DeleteEssay(ctx, essayID)

		// Verify error propagation
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}
