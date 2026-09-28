package video

import (
	"context"

	"inspirate-consulting/internal/models"
)

// ListVideos retrieves a list of all videos from the S3 bucket
func (h *Handler) ListVideos(ctx context.Context) ([]models.Video, error) {

	videos, err := h.VideoRepository.ListVideos(ctx)
	if err != nil {
		return nil, err
	}

	return videos, nil
}
