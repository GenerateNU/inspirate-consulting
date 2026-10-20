package essayreview

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

func (h *Handler) GetEssayReviewStatus(ctx context.Context, essayID uuid.UUID) (*models.EssayReviewStatus, error) {
	status, err := h.EssayReviewRepository.FindEssayReviewStatus(ctx, h.EssayReviewRepository.DB(), essayID)
	if err != nil {
		return nil, err
	}
	if err := h.Authorizer.CanAccessStudent(ctx, status.StudentID); err != nil {
		return nil, err
	}
	return status, nil
}
