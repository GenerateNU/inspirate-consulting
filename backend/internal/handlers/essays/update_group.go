package essays

import (
	"context"

	"inspirate-consulting/internal/models"
)

// UpdateEssayGroup adds an essay to a group, moves it between groups, or removes
// it from its group when essay_group_id is null.
func (h *Handler) UpdateEssayGroup(ctx context.Context, input *models.UpdateEssayGroupInput) (*models.UpdateEssayGroupOutput, error) {
	updatedEssay, err := h.EssayRepository.UpdateEssayGroup(ctx, input.EssayID, input.Body.EssayGroupID)
	if err != nil {
		return nil, err
	}

	return &models.UpdateEssayGroupOutput{
		Body: models.EssayBody{Essay: updatedEssay},
	}, nil
}
