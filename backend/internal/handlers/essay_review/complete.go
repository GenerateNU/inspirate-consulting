package essayreview

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

func (h *Handler) CompleteEssayReview(ctx context.Context, id uuid.UUID) (*models.EssayReviewTransaction, error) {
	// TODO: guard by role - counselors only, and only the requester's assigned
	// counselor.

	completedReview, err := h.EssayReviewRepository.CompleteEssayReview(ctx, id)
	if err != nil {
		return nil, err
	}

	return completedReview, nil
}
