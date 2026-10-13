package user

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
)

func (h *Handler) CreateUser(ctx context.Context, input *models.CreateUserInput, supabase supabase.SupabaseInterface) (*models.CreateUserOutput, error) {
	signup_response, err := supabase.Signup(input.Body.Email, input.Body.Password, auth.Client)
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
