package globalcollege

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) GetGlobalCollege(ctx context.Context, id int64) (*models.GlobalCollege, error) {
	globalCollege, err := h.GlobalCollegeRepository.GetGlobalCollege(ctx, id)
	if err != nil {
		return nil, err
	}
	return globalCollege, nil
}
