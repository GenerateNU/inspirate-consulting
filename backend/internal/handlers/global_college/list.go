package globalcollege

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) ListGlobalColleges(ctx context.Context, input *models.ListGlobalCollegesInput) (*models.ListGlobalCollegesOutput, error) {
	globalColleges, err := h.GlobalCollegeRepository.ListGlobalColleges(ctx)
	if err != nil {
		return nil, err
	}
	return globalColleges, nil
}
