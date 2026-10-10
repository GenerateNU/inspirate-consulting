package essays

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) GetEssaysByGroup(ctx context.Context, input *models.GetEssaysByGroupInput) (*models.GetEssaysByGroupOutput, error) {
	essays, err := h.EssayRepository.GetEssaysByGroup(ctx, input.EssayGroupID)
	if err != nil {
		return nil, err
	}

	return &models.GetEssaysByGroupOutput{
		Body: models.EssayListBody{Essays: essays},
	}, nil
}
