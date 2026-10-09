package media

import (
	"context"
	"inspirate-consulting/internal/models"
)

func (h *Handler) GetMedia(ctx context.Context, id string) (*models.Media, error) {

	fetchedMedia, err := h.MediaRepository.GetMedia(ctx, id)
	if err != nil {
		return nil, err
	}
	return fetchedMedia, nil
}
