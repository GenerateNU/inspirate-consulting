package globalcollege

import (
	"context"

	"inspirate-consulting/internal/models"
)

// CreateGlobalCollege creates a new global college entry in the database
func (h *Handler) CreateGlobalCollege(ctx context.Context, input models.CreateGlobalCollegeRequestBody) (*models.GlobalCollege, error) {

	// allow for manual override of near duplicate check
	if !input.ConfirmDuplicate {
		if err := h.checkNearDuplicate(ctx, input.SchoolName, input.SchoolLocation, nil); err != nil {
			return nil, err
		}
	}

	createdGlobalCollege, err := h.GlobalCollegeRepository.CreateGlobalCollege(ctx, input)
	if err != nil {
		return nil, err
	}

	return createdGlobalCollege, nil
}
