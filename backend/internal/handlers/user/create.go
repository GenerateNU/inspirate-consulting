package user

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
)

func (h *Handler) CreateUser(ctx context.Context, input *models.CreateUserInput, supabase supabase.SupabaseInterface) (*models.CreateUserOutput, error) {
	tempPassword, err := generateTempPassword()
	if err != nil {
		return nil, err
	}

	signup_response, err := supabase.Signup(input.Body.Email, tempPassword, auth.Client)
	if err != nil {
		return nil, err
	}

	user, err := h.UserRepository.CreateUser(ctx, *input, signup_response.User.ID)
	if err != nil {
		return nil, err
	}

	user.Body.TempPassword = tempPassword
	return user, nil
}
