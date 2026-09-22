package essayreview

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

func (h *Handler) RefundEssayReview(ctx context.Context, id uuid.UUID) (*models.RefundEssayReviewResponseBody, error) {
	// TODO: guard by role - students only, and only this transaction's requester.

	refundedReview, err := h.EssayReviewRepository.RefundEssayReview(ctx, id)
	if err != nil {
		return nil, err
	}

	return refundedReview, nil
}
