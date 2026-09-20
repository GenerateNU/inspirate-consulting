package globalcollege

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) GetGlobalCollege(ctx context.Context, input *models.GetGlobalCollegeInput) (*models.GetGlobalCollegeOutput, error) {
	globalCollege, err := h.GlobalCollegeRepository.GetGlobalCollege(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return globalCollege, nil
}
