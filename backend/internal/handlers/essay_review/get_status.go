package essayreview

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

func (h *Handler) GetEssayReviewStatus(ctx context.Context, essayID uuid.UUID) (*models.EssayReviewStatus, error) {
	// TODO: guard by role - the essay's owner or their assigned counselor.

	return h.EssayReviewRepository.FindEssayReviewStatus(ctx, h.EssayReviewRepository.DB(), essayID)
}
