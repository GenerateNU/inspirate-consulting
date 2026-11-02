package user

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// updates the logged-in user's name and pfp_key
func (h *Handler) UpdateUser(ctx context.Context, input models.UpdateUserBody) (*models.User, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	updated, err := h.UserRepository.UpdateUser(ctx, userID, input)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
