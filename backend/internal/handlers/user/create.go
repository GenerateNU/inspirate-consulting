package user

import (
	"context"
	"fmt"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/supabase"
	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateUser(ctx context.Context, input *models.CreateUserInput, supabase supabase.SupabaseInterface) (*models.CreateUserOutput, error) {
	signup_response, err := supabase.Signup(input.Body.Email, input.Body.Password, auth.Client)
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
