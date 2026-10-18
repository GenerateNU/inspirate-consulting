package counselor

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
)

// Unit tests for the ListPendingReviews handler logic
// Written with stub counselor ID
func TestHandler_ListPendingReviews(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	counselorID := auth.GetCounselorID(ctx)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expectedReviews := []models.PendingReview{
			{TransactionID: uuid.New(), EssayID: uuid.New(), EssayType: "personal-statement", EssayLink: "https://example.com/a", StudentID: uuid.New(), StudentName: "Alice"},
			{TransactionID: uuid.New(), EssayID: uuid.New(), EssayType: "supplemental", EssayLink: "https://example.com/b", StudentID: uuid.New(), StudentName: "Bob"},
		}

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListPendingReviews", ctx, counselorID).Return(expectedReviews, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.ListPendingReviews(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedReviews, res)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListPendingReviews", ctx, counselorID).Return([]models.PendingReview{}, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.ListPendingReviews(ctx)

		assert.NoError(t, err)
		assert.Equal(t, []models.PendingReview{}, res)
	})

	t.Run("repository error propagates", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListPendingReviews", ctx, counselorID).Return(nil, errors.New("database error"))

		handler := NewHandler(mockRepo)
		res, err := handler.ListPendingReviews(ctx)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "database error")
	})
}
