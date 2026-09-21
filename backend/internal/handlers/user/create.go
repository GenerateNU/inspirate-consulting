package user

import (
	"context"
	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateUser(ctx context.Context, input *models.CreateUserInput) (*models.CreateUserOutput, error) {
	user, err := h.UserRepository.CreateUser(ctx, *input)
	if err != nil {
		return nil, err
	}
	// fire off supabase acc creation method

	return user, nil
}
