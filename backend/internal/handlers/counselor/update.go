package counselor

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// updates the logged-in counselor's profile fields
func (h *Handler) UpdateCounselor(ctx context.Context, input models.UpdateCounselorBody) (*models.Counselor, error) {
	counselorID, err := auth.GetCounselorID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	updated, err := h.CounselorRepository.UpdateCounselor(ctx, counselorID, input)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
