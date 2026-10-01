package video

import (
	"context"

	"inspirate-consulting/internal/models"
)

// GetVideo retrieves a video from the S3 bucket
func (h *Handler) GetVideo(ctx context.Context, s3Key string) (*models.Video, error) {

	video, err := h.VideoRepository.GetVideo(ctx, s3Key)
	if err != nil {
		return nil, err
	}

	return video, nil
}
