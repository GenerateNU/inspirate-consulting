package student

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) SetReviewBalance(ctx context.Context, id uuid.UUID, input models.SetStudentReviewBalanceRequestBody) (*models.Student, error) {
	// TODO: guard by role - admins and counselors only.

	if input.ReviewBalance < 0 {
		return nil, errs.BadRequest("review_balance cannot be negative")
	}

	updatedStudent, err := h.EssayReviewRepository.SetReviewBalance(ctx, id, input.ReviewBalance)
	if err != nil {
		return nil, err
	}

	return updatedStudent, nil
}
