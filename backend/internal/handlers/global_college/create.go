package globalcollege

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateGlobalCollege(ctx context.Context, input *models.CreateGlobalCollegeInput) (*models.CreateGlobalCollegeOutput, error) {
	globalCollege, err := h.GlobalCollegeRepository.CreateGlobalCollege(ctx, *input)
	if err != nil {
		return nil, err
	}
	return globalCollege, nil
}
