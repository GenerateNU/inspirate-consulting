package media

import (
	"context"
	"inspirate-consulting/internal/models"
)

func (h *Handler) UpdateMedia(ctx context.Context, id string, change string) (*models.Media, error) {
	updatedMedia, err := h.MediaRepository.UpdateMedia(ctx, id, change)

	if err != nil {
		return nil, err
	}

	return updatedMedia, err
}
