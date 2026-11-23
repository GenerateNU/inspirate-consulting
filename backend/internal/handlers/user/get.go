package user

import (
	"context"
	"inspirate-consulting/internal/models"
	"log/slog"
)

func (h *Handler) FetchUser(ctx context.Context, input *models.FetchUserInput) (*models.FetchUserProfileOutput, error) {
	user, err := h.UserRepository.FetchUser(ctx, *input)
	if err != nil {
		return nil, err
	}

	profile := h.withPfpURL(ctx, user.Body)
	return &models.FetchUserProfileOutput{Body: &profile}, nil
}

// withPfpURL presigns a download link for the user's picture. A broken picture
// link shouldn't block loading the user, so failures only drop the URL.
func (h *Handler) withPfpURL(ctx context.Context, user *models.User) models.UserProfile {
	profile := models.UserProfile{User: *user}
	if user.PfpKey == nil {
		return profile
	}
	if url, err := h.PfpRepository.GetDownloadURL(ctx, *user.PfpKey); err == nil {
		profile.PfpURL = &url
	} else {
		slog.Error("failed to presign profile picture", "user_id", user.ID, "err", err)
	}
	return profile
}
