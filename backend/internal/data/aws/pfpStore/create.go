package pfpRepository

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *PfpRepository) PresignUpload(ctx context.Context, userID uuid.UUID, input models.PresignPfpUploadRequestBody) (*models.PresignUploadResponse, error) {
	ext, ok := pfpExtensions[input.ContentType]
	if !ok {
		return nil, errs.BadRequest("profile picture must be a JPEG, PNG, or WebP image")
	}
	if input.ContentLength < 1 || input.ContentLength > maxPfpBytes {
		return nil, errs.BadRequest("profile picture must be at most 5 MiB")
	}

	key := userPrefix(userID) + uuid.NewString() + ext

	uploadURL, err := r.s3.PresignUpload(ctx, key, input.ContentType, input.ContentLength)
	if err != nil {
		return nil, err
	}

	return &models.PresignUploadResponse{
		S3Key:     key,
		UploadURL: uploadURL,
	}, nil
}
