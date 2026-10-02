package essays

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) UpdateStatus(ctx context.Context, input *models.UpdateStatusInput) (*models.UpdateStatusOutput, error) {
	updatedEssay, err := h.EssayRepository.UpdateStatus(ctx, input.EssayID, input.Body.Status)
	if err != nil {
		return nil, err
	}

	return &models.UpdateStatusOutput{
		Essay: updatedEssay,
	}, nil
}
