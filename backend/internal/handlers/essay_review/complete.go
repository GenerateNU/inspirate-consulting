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

func (h *Handler) CompleteEssayReview(ctx context.Context, id uuid.UUID) (*models.EssayReviewTransaction, error) {
	// TODO: guard by role - counselors only, and only the requester's assigned
	// counselor.

	var completed *models.EssayReviewTransaction

	err := h.EssayReviewRepository.WithTx(ctx, func(db dbinterface.QueryInterface) error {
		review, err := h.EssayReviewRepository.LockTransaction(ctx, db, id)
		if err != nil {
			return err
		}

		if err := completable(review); err != nil {
			return err
		}

		completed, err = h.EssayReviewRepository.MarkCompleted(ctx, db, id)
		return err
	})
	if err != nil {
		return nil, err
	}

	return completed, nil
}

func completable(review *models.EssayReviewTransaction) error {
	if review.Status == nil {
		return errs.BadRequest(fmt.Sprintf(
			"only a spend can be completed, this entry_type is %s", review.EntryType,
		))
	}

	switch *review.Status {
	case models.ReviewStatusRefunded:
		return errs.NewHTTPError(http.StatusConflict, errors.New(
			"review has already been refunded and cannot be completed",
		))
	case models.ReviewStatusCompleted:
		return errs.NewHTTPError(http.StatusConflict, errors.New(
			"review has already been completed",
		))
	}

	return nil
}
