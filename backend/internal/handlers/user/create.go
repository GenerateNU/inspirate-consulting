package user

import (
	"context"
	"fmt"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateUser(ctx context.Context, input *models.CreateUserInput, config *config.Config) (*models.CreateUserOutput, error) {
	fmt.Println("[handler] CreateUser called")

	signup_response, err := auth.SupabaseSignup(&config.Supabase, input.Body.Email, input.Body.Password)
	if err != nil {
		fmt.Println("[handler] SupabaseSignup error:", err)
		return nil, err
	}
	fmt.Println("[handler] SupabaseSignup success")

	supabase_id := signup_response.User.ID
	user, err := h.UserRepository.CreateUser(ctx, *input, supabase_id)
	if err != nil {
		fmt.Println("[handler] DB CreateUser error:", err)
		return nil, err
	}

	return user, nil
}
