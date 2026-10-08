package media

import(
	"context"
	"inspirate-consulting/internal/models"
)

func(h *Handler) ListAllMedia(ctx context.Context, limit int, offset int)([]models.Media, error){
	allMedia, err := h.MediaRepository.ListAllMedia(ctx, limit, offset)
	if(err != nil) {
		return nil, err
	}
	return allMedia, nil
}