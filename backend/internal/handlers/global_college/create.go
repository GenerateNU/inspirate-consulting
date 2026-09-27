package globalcollege

import (
	"context"

	"inspirate-consulting/internal/models"
)

// CreateGlobalCollege creates a new global college entry in the database
func (h *Handler) CreateGlobalCollege(ctx context.Context, input models.CreateGlobalCollegeRequestBody) (*models.GlobalCollege, error) {

	createdGlobalCollege, err := h.GlobalCollegeRepository.CreateGlobalCollege(ctx, input)
	if err != nil {
		return nil, err
	}
 
	return createdGlobalCollege, nil
}
