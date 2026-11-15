package pfp

import (
	"context"
	"log/slog"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
)

// clears the logged-in user's profile picture and removes the image from S3
func (h *Handler) DeletePfp(ctx context.Context) error {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return errs.Unauthorized()
	}

	_, previousKey, err := h.UserRepository.SetPfpKey(ctx, userID, nil)
	if err != nil {
		return err
	}
	if previousKey == nil {
		return errs.NotFound("profile picture", "user id", userID.String())
	}

	// the user no longer points at the image, so a failed cleanup only leaves an orphaned object
	if err := h.PfpRepository.DeletePfp(ctx, userID, *previousKey); err != nil {
		slog.Error("failed to delete profile picture", "user_id", userID, "s3_key", *previousKey, "err", err)
	}
	return nil
}
