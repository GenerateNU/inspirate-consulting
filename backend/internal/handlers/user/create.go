package user

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateUser(ctx context.Context, input *models.CreateUserInput, config *config.Config) (*models.CreateUserOutput, error) {

	// fire off supabase acc creation method

	signup_response, err := auth.SupabaseSignup(&config.Supabase, input.Body.Email, input.Body.Password)
	if err != nil {
		return nil, err
	}

	supabase_id := signup_response.User.ID
	user, err := h.UserRepository.CreateUser(ctx, *input, supabase_id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
