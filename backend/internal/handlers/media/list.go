package media

import (
	"context"
	"inspirate-consulting/internal/models"
)

// minimum pg_trgm similarity, from 0.0 (complete mismatch) to 1.0 (perfect match), for a media's
// title or description to count as matching the search query
const searchThreshold = 0.2

func (h *Handler) ListAllMedia(ctx context.Context, query string, limit int, offset int) ([]models.Media, error) {
	allMedia, err := h.MediaRepository.ListAllMedia(ctx, query, searchThreshold, limit, offset)
	if err != nil {
		return nil, err
	}
	return allMedia, nil
}
