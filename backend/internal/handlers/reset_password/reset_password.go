package resetpassword

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
)

func (h *Handler) ResetPassword(ctx context.Context, input *models.ResetPasswordInput, supabase supabase.SupabaseInterface) (*models.ResetPasswordResponse, error) {

	// also needs to update field in user table, goes here

	res, err := supabase.SupabaseResetPassword(auth.Client, input.Body.NewPassword, input.UserID)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
