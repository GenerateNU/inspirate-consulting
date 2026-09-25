package essayreview

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// Unit tests for the CompleteEssayReview handler logic without HTTP or database dependencies.
func TestHandler_CompleteEssayReview(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	rejections := []struct {
		name       string
		review     *models.EssayReviewTransaction
		wantStatus int
	}{
		{"a reversal row is not a review", transactionRow(func(row *models.EssayReviewTransaction) {
			row.EntryType = models.ReviewEntryRefund
			row.Subtotal = 1
		}), 400},
		{"an adjustment row is not a review", transactionRow(func(row *models.EssayReviewTransaction) {
			row.EntryType = models.ReviewEntryAdjustment
			row.Subtotal = 5
			row.EssayID = nil
		}), 400},
		{"an already refunded review", transactionRow(func(row *models.EssayReviewTransaction) {
			reversalID := uuid.New()
			row.Refund = &reversalID
		}), 409},
		{"an already completed review", transactionRow(func(row *models.EssayReviewTransaction) {
			row.CompletedAt = &fixedTime
		}), 409},
	}

	for _, tc := range rejections {
		t.Run(tc.name+" is rejected without writing", func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewEssayReviewRepository(t)
			expectTx(mockRepo)
			mockRepo.On("LockTransaction", mock.Anything, mock.Anything, tc.review.ID).Return(tc.review, nil)

			handler := NewHandler(mockRepo)
			res, err := handler.CompleteEssayReview(ctx, tc.review.ID)

			assert.Nil(t, res)
			require.Error(t, err)
			assert.Equal(t, tc.wantStatus, err.(errs.HTTPError).Code)
			mockRepo.AssertNotCalled(t, "MarkCompleted", mock.Anything, mock.Anything, mock.Anything)
		})
	}

	t.Run("success marks the review completed and leaves the balance alone", func(t *testing.T) {
		t.Parallel()

		studentID, essayID := uuid.New(), uuid.New()
		review := transactionRow(func(row *models.EssayReviewTransaction) {
			row.StudentID = studentID
			row.EssayID = &essayID
		})

		completedStatus := models.ReviewStatusCompleted
		completed := &models.EssayReviewTransaction{
			ID: review.ID, StudentID: studentID, EssayID: &essayID, Subtotal: -1,
			EntryType: models.ReviewEntrySpend, Status: &completedStatus,
		}

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockTransaction", mock.Anything, mock.Anything, review.ID).Return(review, nil)
		mockRepo.On("MarkCompleted", mock.Anything, mock.Anything, review.ID).Return(completed, nil)

		handler := NewHandler(mockRepo)
		res, err := handler.CompleteEssayReview(ctx, review.ID)

		require.NoError(t, err)
		assert.Equal(t, completed, res)
		// completing spends nothing: the charge was taken at request time
		mockRepo.AssertNotCalled(t, "AdjustStudentBalance", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("unknown review propagates the repository error", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockTransaction", mock.Anything, mock.Anything, id).
			Return(nil, errs.NotFound("essay review", "id", id.String()))

		handler := NewHandler(mockRepo)
		res, err := handler.CompleteEssayReview(ctx, id)

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 404, err.(errs.HTTPError).Code)
	})
}
