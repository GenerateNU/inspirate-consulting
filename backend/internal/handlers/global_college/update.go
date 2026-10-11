package globalcollege

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) UpdateGlobalCollege(ctx context.Context, id int64, input models.UpdateGlobalCollegeRequestBody) (*models.GlobalCollege, error) {
	updated, err := h.GlobalCollegeRepository.UpdateGlobalCollege(ctx, id, input, nearDuplicateThreshold)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
