package globalcollege

import (
	"context"

	"inspirate-consulting/internal/models"
)

// range from 0.0 to 1.0, inclusive, where 0.0 is a complete mismatch and 1.0 is a perfect match
const nearDuplicateThreshold = 0.45

// CreateGlobalCollege creates a new global college entry in the database
func (h *Handler) CreateGlobalCollege(ctx context.Context, input models.CreateGlobalCollegeRequestBody) (*models.GlobalCollege, error) {

	createdGlobalCollege, err := h.GlobalCollegeRepository.CreateGlobalCollege(ctx, input, nearDuplicateThreshold)
	if err != nil {
		return nil, err
	}

	return createdGlobalCollege, nil
}
