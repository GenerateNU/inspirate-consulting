package videoRepository
 
import (
	"context"
 
	"inspirate-consulting/internal/data/aws/s3client"
	"inspirate-consulting/internal/models"
)

func (r *VideoRepository) PresignUpload(ctx context.Context, originalFilename string) (*models.PresignUploadResponse, error) {
	key := s3client.GenerateObjectKey(originalFilename)

	uploadURL, err := r.s3.PresignUpload(ctx, key)
	if err != nil {
		return nil, err
	}

	return &models.PresignUploadResponse{
		S3Key:     key,
		UploadURL: uploadURL,
	}, nil
}