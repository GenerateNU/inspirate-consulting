package pfpRepository

import (
	"context"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
)

func (r *PfpRepository) DeletePfp(ctx context.Context, userID uuid.UUID, s3Key string) error {
	if !isUserPfpKey(userID, s3Key) {
		return errs.Forbidden("you can only delete your own profile picture")
	}
	return r.s3.DeleteObject(ctx, s3Key)
}
