package essay_groups

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateEssayGroup(ctx context.Context, input *models.CreateEssayGroupInput) (*models.CreateEssayGroupOutput, error) {
	createdGroup, err := h.EssayGroupRepository.CreateEssayGroup(ctx, input.Body)
	if err != nil {
		return nil, err
	}

	return &models.CreateEssayGroupOutput{
		Body: *createdGroup,
	}, nil
}
