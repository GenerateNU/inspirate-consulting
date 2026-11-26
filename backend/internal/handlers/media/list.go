package media

import (
	"context"
	"inspirate-consulting/internal/models"
)

// minimum pg_trgm word_similarity, from 0.0 (complete mismatch) to 1.0 (perfect match), between the
// search query and the best-matching part of a media's title or description
const searchThreshold = 0.4

func (h *Handler) ListAllMedia(ctx context.Context, query string, limit int, offset int) ([]models.Media, error) {
	allMedia, err := h.MediaRepository.ListAllMedia(ctx, query, searchThreshold, limit, offset)
	if err != nil {
		return nil, err
	}
	return allMedia, nil
}
