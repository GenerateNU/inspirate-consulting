package essayreview

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) RefundEssayReview(ctx context.Context, id uuid.UUID) (*models.RefundEssayReviewResponseBody, error) {
	// TODO: guard by role - students only, and only this transaction's requester.

	result := &models.RefundEssayReviewResponseBody{}

	err := h.EssayReviewRepository.WithTx(ctx, func(db dbinterface.QueryInterface) error {
		charge, err := h.EssayReviewRepository.LockTransaction(ctx, db, id)
		if err != nil {
			return err
		}

		if err := refundable(charge); err != nil {
			return err
		}

		if _, err := h.EssayReviewRepository.LockStudent(ctx, db, charge.StudentID); err != nil {
			return err
		}

		reversal, err := h.EssayReviewRepository.InsertRefund(ctx, db, charge)
		if err != nil {
			return err
		}

		refunded, err := h.EssayReviewRepository.LinkRefund(ctx, db, charge.ID, reversal.ID)
		if err != nil {
			return err
		}

		if err := h.EssayReviewRepository.AdjustStudentBalance(ctx, db, charge.StudentID, reversal.Subtotal); err != nil {
			return err
		}

		result.Transaction = *refunded
		result.Refund = *reversal
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func refundable(charge *models.EssayReviewTransaction) error {
	if charge.Status == nil {
		return errs.BadRequest(fmt.Sprintf(
			"only a spend can be refunded, this entry_type is %s", charge.EntryType,
		))
	}

	switch *charge.Status {
	case models.ReviewStatusCompleted:
		return errs.NewHTTPError(http.StatusConflict, errors.New(
			"review has already been completed and cannot be refunded",
		))
	case models.ReviewStatusRefunded:
		return errs.NewHTTPError(http.StatusConflict, errors.New(
			"review has already been refunded",
		))
	}

	return nil
}
