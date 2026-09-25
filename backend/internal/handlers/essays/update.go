package essays

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) UpdateStatus(ctx context.Context, input *models.UpdateStatusInput) (*models.UpdateStatusOutput, error) {
	if err := h.EssayRepository.UpdateStatus(ctx, input.EssayID, input.Body.Status); err != nil {
		return nil, err
	}

	return &models.UpdateStatusOutput{}, nil
}
