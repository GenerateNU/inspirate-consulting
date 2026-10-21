package counselor

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// lists the open essay reviews the authenticated counselor needs to complete, oldest first
func (h *Handler) ListPendingReviews(ctx context.Context) ([]models.PendingReview, error) {

	counselorID, err := auth.GetCounselorID(ctx)
	if err != nil {
		return nil, err
	}

	reviews, err := h.CounselorRepository.ListPendingReviews(ctx, counselorID)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
