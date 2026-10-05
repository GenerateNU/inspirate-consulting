package media

import(
	"context"
	"inspirate-consulting/internal/models"
)

func(h *Handler) ListAllMedia(ctx context.Context)([]models.Media, error){
	allMedia, err := h.MediaRepository.ListAllMedia(ctx)
	if(err != nil) {
		return nil, err
	}
	return allMedia, nil
}