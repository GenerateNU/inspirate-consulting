package video

import (
	"context"

	"inspirate-consulting/internal/models"
)

// UploadVideo uploads a new video to the S3 bucket
func (h *Handler) UploadVideo(ctx context.Context, input models.PresignUploadRequestBody) (*models.PresignUploadResponse, error) {

	presignUploadResponse, err := h.VideoRepository.PresignUpload(ctx, input)
	if err != nil {
		return nil, err
	}

	return presignUploadResponse, nil
}
