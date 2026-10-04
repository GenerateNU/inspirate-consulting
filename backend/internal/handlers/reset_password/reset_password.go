package resetpassword

import (
	"context"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
)

func (h *Handler) ResetPassword(ctx context.Context, input *models.ResetPasswordInput, supabase supabase.SupabaseInterface) (*models.ResetPasswordResponse, error) {
	return nil, nil
}
