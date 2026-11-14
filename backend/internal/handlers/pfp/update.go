package pfp

import (
	"context"
	"log/slog"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// saves an uploaded image as the logged-in user's profile picture and removes the one it replaces
func (h *Handler) UpdatePfp(ctx context.Context, s3Key string) (*models.Pfp, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	if err := h.PfpRepository.VerifyUpload(ctx, userID, s3Key); err != nil {
		return nil, err
	}

	_, previousKey, err := h.UserRepository.SetPfpKey(ctx, userID, &s3Key)
	if err != nil {
		return nil, err
	}

	// the new picture is already saved, so a failed cleanup only leaves an orphaned object
	if previousKey != nil && *previousKey != s3Key {
		if err := h.PfpRepository.DeletePfp(ctx, userID, *previousKey); err != nil {
			slog.Error("failed to delete replaced profile picture", "user_id", userID, "s3_key", *previousKey, "err", err)
		}
	}

	downloadURL, err := h.PfpRepository.GetDownloadURL(ctx, s3Key)
	if err != nil {
		return nil, err
	}

	return &models.Pfp{S3Key: s3Key, DownloadURL: downloadURL}, nil
}
