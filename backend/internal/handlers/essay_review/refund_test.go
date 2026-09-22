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

// Unit tests for the RefundEssayReview handler logic without HTTP or database dependencies.
func TestHandler_RefundEssayReview(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	rejections := []struct {
		name       string
		charge     *models.EssayReviewTransaction
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
		{"an already completed review", transactionRow(func(row *models.EssayReviewTransaction) {
			row.CompletedAt = &fixedTime
		}), 409},
		{"an already refunded review", transactionRow(func(row *models.EssayReviewTransaction) {
			reversalID := uuid.New()
			row.Refund = &reversalID
		}), 409},
	}

	for _, tc := range rejections {
		t.Run(tc.name+" is rejected without writing", func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewEssayReviewRepository(t)
			expectTx(mockRepo)
			mockRepo.On("LockTransaction", mock.Anything, mock.Anything, tc.charge.ID).Return(tc.charge, nil)

			handler := NewHandler(mockRepo)
			res, err := handler.RefundEssayReview(ctx, tc.charge.ID)

			assert.Nil(t, res)
			require.Error(t, err)
			assert.Equal(t, tc.wantStatus, err.(errs.HTTPError).Code)
			mockRepo.AssertNotCalled(t, "InsertRefund", mock.Anything, mock.Anything, mock.Anything)
			mockRepo.AssertNotCalled(t, "LinkRefund", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			mockRepo.AssertNotCalled(t, "AdjustStudentBalance", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		})
	}

	t.Run("success inserts a reversal, links it, then credits the balance", func(t *testing.T) {
		t.Parallel()

		studentID, essayID := uuid.New(), uuid.New()
		charge := transactionRow(func(row *models.EssayReviewTransaction) {
			row.StudentID = studentID
			row.EssayID = &essayID
		})

		reversal := &models.EssayReviewTransaction{
			ID: uuid.New(), StudentID: studentID, EssayID: &essayID,
			Subtotal: 1, EntryType: models.ReviewEntryRefund,
		}
		refundedStatus := models.ReviewStatusRefunded
		refunded := &models.EssayReviewTransaction{
			ID: charge.ID, StudentID: studentID, EssayID: &essayID, Subtotal: -1,
			EntryType: models.ReviewEntrySpend, Status: &refundedStatus, Refund: &reversal.ID,
		}

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockTransaction", mock.Anything, mock.Anything, charge.ID).Return(charge, nil)
		mockRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(&models.Student{ID: studentID}, nil)
		mockRepo.On("InsertRefund", mock.Anything, mock.Anything, charge).Return(reversal, nil)
		mockRepo.On("LinkRefund", mock.Anything, mock.Anything, charge.ID, reversal.ID).Return(refunded, nil)
		mockRepo.On("AdjustStudentBalance", mock.Anything, mock.Anything, studentID, 1).Return(nil)

		handler := NewHandler(mockRepo)
		res, err := handler.RefundEssayReview(ctx, charge.ID)

		require.NoError(t, err)
		assert.Equal(t, *refunded, res.Transaction)
		assert.Equal(t, *reversal, res.Refund)
	})

	t.Run("unknown review propagates the repository error", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockTransaction", mock.Anything, mock.Anything, id).
			Return(nil, errs.NotFound("essay review", "id", id.String()))

		handler := NewHandler(mockRepo)
		res, err := handler.RefundEssayReview(ctx, id)

		assert.Nil(t, res)
		require.Error(t, err)
		assert.Equal(t, 404, err.(errs.HTTPError).Code)
	})
}
