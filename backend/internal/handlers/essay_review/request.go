package essayreview

import (
	"context"
	"fmt"
	"net/http"

	"inspirate-consulting/internal/auth"
	dbinterface "inspirate-consulting/internal/data/db-interface"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) RequestEssayReview(ctx context.Context, input models.RequestEssayReviewRequestBody) (*models.EssayReviewTransaction, error) {
	studentID, err := auth.GetStudentID(ctx)
	if err != nil {
		return nil, err
	}

	if input.Amount < 1 {
		return nil, errs.BadRequest("amount must be at least 1")
	}

	essayOwnerID, err := h.EssayReviewRepository.FindEssayStudentID(ctx, h.EssayReviewRepository.DB(), input.EssayID)
	if err != nil {
		return nil, err
	}
	if essayOwnerID != studentID {
		return nil, errs.Forbidden("you can only request reviews for your own essays")
	}

	var spend *models.EssayReviewTransaction

	err = h.EssayReviewRepository.WithTx(ctx, func(db dbinterface.QueryInterface) error {

		student, err := h.EssayReviewRepository.LockStudent(ctx, db, studentID)
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

		if err := h.EssayReviewRepository.AdjustStudentBalance(ctx, db, studentID, -input.Amount); err != nil {
			return err
		}

		spend, err = h.EssayReviewRepository.InsertSpend(ctx, db, studentID, input.EssayID, input.Amount)
		return err
	})
	if err != nil {
		return nil, err
	}

	return spend, nil
}
