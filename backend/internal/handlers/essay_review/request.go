package essayreview

import (
	"context"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) RequestEssayReview(ctx context.Context, input models.RequestEssayReviewRequestBody) (*models.EssayReviewTransaction, error) {
	// TODO: guard by role - students only, input.StudentID must be the caller,
	// and the essay must belong to them.

	if input.Amount < 1 {
		return nil, errs.BadRequest("amount must be at least 1")
	}

	requestedReview, err := h.EssayReviewRepository.RequestEssayReview(ctx, input)
	if err != nil {
		return nil, err
	}

	return requestedReview, nil
}
