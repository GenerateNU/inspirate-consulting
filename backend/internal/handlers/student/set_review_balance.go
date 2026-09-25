package student

import (
	"context"

	"github.com/google/uuid"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) SetReviewBalance(ctx context.Context, id uuid.UUID, input models.SetStudentReviewBalanceRequestBody) (*models.Student, error) {
	// TODO: guard by role - admins and counselors only. This mints review credit
	// from nothing, so a student must never reach it.

	if input.ReviewBalance < 0 {
		return nil, errs.BadRequest("review_balance cannot be negative")
	}

	var updated *models.Student

	err := h.EssayReviewRepository.WithTx(ctx, func(db dbinterface.QueryInterface) error {
		student, err := h.EssayReviewRepository.LockStudent(ctx, db, id)
		if err != nil {
			return err
		}

		delta := input.ReviewBalance - student.ReviewBalance
		if delta == 0 {
			updated = student
			return nil
		}

		if err := h.EssayReviewRepository.SetStudentBalance(ctx, db, id, input.ReviewBalance); err != nil {
			return err
		}
		if err := h.EssayReviewRepository.InsertAdjustment(ctx, db, id, delta); err != nil {
			return err
		}

		student.ReviewBalance = input.ReviewBalance
		updated = student
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updated, nil
}
