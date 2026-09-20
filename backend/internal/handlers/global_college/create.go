package globalcollege

import (
	"context"

	"inspirate-consulting/internal/models"
)

// CreateGlobalCollege creates a new global college entry in the database
func (h *Handler) CreateGlobalCollege(ctx context.Context, input *models.CreateGlobalCollegeInput) (*models.CreateGlobalCollegeOutput, error) {

	createdGlobalCollege, err := h.GlobalCollegeRepository.CreateGlobalCollege(ctx, input.Body)
	if err != nil {
		return nil, err
	}
 
	return &models.CreateGlobalCollegeOutput{Body: *createdGlobalCollege}, nil
}
