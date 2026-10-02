package user

import (
	"context"
	"inspirate-consulting/internal/models"
)

func (h *Handler) FetchUser(ctx context.Context, input *models.FetchUserInput) (*models.FetchUserOutput, error) {

	user, err := h.UserRepository.FetchUser(ctx, *input)
	if err != nil {
		return nil, err
	}

	return user, nil
}
