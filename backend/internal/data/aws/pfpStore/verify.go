package pfpRepository

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
)

// checks that the key is the user's own and that S3 holds an allowed image there
func (r *PfpRepository) VerifyUpload(ctx context.Context, userID uuid.UUID, s3Key string) error {
	if !isUserPfpKey(userID, s3Key) {
		return errs.Forbidden("you can only use your own profile picture")
	}

	metadata, err := r.s3.GetObjectMetadata(ctx, s3Key)
	if err != nil {
		return err
	}
	if metadata == nil {
		return errs.NotFound("profile picture", "s3Key", s3Key)
	}

	if _, ok := pfpExtensions[metadata.ContentType]; !ok || metadata.ContentLength > maxPfpBytes {
		return errs.BadRequest("profile picture must be a JPEG, PNG, or WebP image of at most 5 MiB")
	}
	return nil
}
