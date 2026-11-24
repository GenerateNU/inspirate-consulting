package videoRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"inspirate-consulting/internal/models"
)

// keeps the original filename in the key for readability
func generateVideoKey(originalFilename string) string {
	return fmt.Sprintf("%s%s-%s", videoPrefix, uuid.NewString(), originalFilename)
}

func (r *VideoRepository) PresignUpload(ctx context.Context, input models.PresignUploadRequestBody) (*models.PresignUploadResponse, error) {
	key := generateVideoKey(input.OriginalFilename)

	uploadURL, err := r.s3.PresignUpload(ctx, key, input.ContentType, input.ContentLength)
	if err != nil {
		return nil, err
	}

	return &models.PresignUploadResponse{
		S3Key:     key,
		UploadURL: uploadURL,
	}, nil
}
