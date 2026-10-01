package videoRepository
 
import (
	"context"
 
	"inspirate-consulting/internal/models"
)

func (r *VideoRepository) ListVideos(ctx context.Context) ([]models.Video, error) {
	// get the object keys from the S3 bucket
	keys, err := r.s3.ListObjectKeys(ctx)
	if err != nil {
		return nil, err
	}

	// fetch the presigned download URLs for each object key
	videos := make([]models.Video, 0, len(keys))
	for _, key := range keys {
		url, err := r.s3.PresignDownload(ctx, key)
		if err != nil {
			return nil, err
		}
		videos = append(videos, models.Video{S3Key: key, DownloadURL: url})
	}

	return videos, nil
}