package login

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
)

func (h *Handler) Login(ctx context.Context, input *models.LoginInput, supabase supabase.SupabaseInterface) (*models.LoginResponse, error) {
	res, err := supabase.SupabaseLogin(input.Body.Email, input.Body.Password, auth.Client)
	if err != nil {
		return nil, err
	}

	userInput := models.FetchUserInput{}
	userInput.ID = res.User.ID
	user, err := h.LoginRepository.FetchUser(ctx, userInput)
	if err != nil {
		return nil, err
	}

	res.ResetTime = *user.Body.ResetTime
	return &res, nil
}
