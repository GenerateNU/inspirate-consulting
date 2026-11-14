package videoRepository

import (
	"context"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *VideoRepository) GetVideo(ctx context.Context, s3Key string) (*models.Video, error) {
	if !isVideoKey(s3Key) {
		return nil, errs.NotFound("video", "s3Key", s3Key)
	}

	url, err := r.s3.PresignDownload(ctx, s3Key)
	if err != nil {
		return nil, err
	}

	return &models.Video{
		S3Key:       s3Key,
		DownloadURL: url,
	}, nil
}
