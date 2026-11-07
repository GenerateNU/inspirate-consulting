package videoRepository

import (
	"context"
	"inspirate-consulting/internal/errs"
)

func (r *VideoRepository) DeleteVideo(ctx context.Context, s3Key string) error {
	exists, err := r.s3.ObjectExists(ctx, s3Key)

	if err != nil {
		return err
	}
	if !exists {
		return errs.NotFound("video", "s3Key", s3Key)
	}

	return r.s3.DeleteObject(ctx, s3Key)

}
