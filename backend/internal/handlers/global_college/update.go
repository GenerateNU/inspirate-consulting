package globalcollege

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) UpdateGlobalCollege(ctx context.Context, id int64, input models.UpdateGlobalCollegeRequestBody) (*models.GlobalCollege, error) {

	// allow for manual override of near duplicate check
	if !input.ConfirmDuplicate {
		if err := h.checkNearDuplicate(ctx, input.SchoolName, input.SchoolLocation, &id); err != nil {
			return nil, err
		}
	}

	updated, err := h.GlobalCollegeRepository.UpdateGlobalCollege(ctx, id, input)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
