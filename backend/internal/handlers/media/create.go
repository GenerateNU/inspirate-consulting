package media

import (
	"context"
	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateMedia(ctx context.Context, item *models.CreateMediaRequestBody) (*models.Media, error) {
	createdMedia, err := h.MediaRepository.CreateMedia(ctx, item)

	if err != nil {
		return nil, err
	}
	return createdMedia, nil

}
