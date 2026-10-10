package essayreview

import (
	"context"
	"fmt"
	"net/http"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) RequestEssayReview(ctx context.Context, input models.RequestEssayReviewRequestBody) (*models.EssayReviewTransaction, error) {
	// TODO: guard by role - students only, input.StudentID must be the caller,
	// and the essay must belong to them.

	if input.Amount < 1 {
		return nil, errs.BadRequest("amount must be at least 1")
	}

	var spend *models.EssayReviewTransaction

	err := h.EssayReviewRepository.WithTx(ctx, func(db dbinterface.QueryInterface) error {

		student, err := h.EssayReviewRepository.LockStudent(ctx, db, input.StudentID)
		if err != nil {
			return err
		}

		openReviewID, err := h.EssayReviewRepository.FindOpenReviewForEssay(ctx, db, input.EssayID)
		if err != nil {
			return err
		}
		if openReviewID != nil {
			return errs.NewHTTPError(http.StatusConflict, fmt.Errorf(
				"essay already has an open review request (%s)", *openReviewID,
			))
		}

		if student.ReviewBalance < input.Amount {
			return errs.BadRequest(fmt.Sprintf(
				"insufficient review balance: student has %d, requested %d", student.ReviewBalance, input.Amount,
			))
		}

		if err := h.EssayReviewRepository.AdjustStudentBalance(ctx, db, input.StudentID, -input.Amount); err != nil {
			return err
		}

		spend, err = h.EssayReviewRepository.InsertSpend(ctx, db, input.StudentID, input.EssayID, input.Amount)
		return err
	})
	if err != nil {
		return nil, err
	}

	return spend, nil
}
