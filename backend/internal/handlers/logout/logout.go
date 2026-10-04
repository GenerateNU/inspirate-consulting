package logout

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
)

func (h *Handler) Logout(ctx context.Context, input *models.LogoutInput, supabase supabase.SupabaseInterface) (*models.LogoutResponse, error) {

	res, err := supabase.SupabaseLogout(auth.Client, input.AccessToken)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
