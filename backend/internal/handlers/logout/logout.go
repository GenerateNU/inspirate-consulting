package logout

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
)

func (h *Handler) Logout(ctx context.Context, input *models.LogoutInput, supabase supabase.SupabaseInterface) (*models.LogoutResponse, error) {
	token, ok := ctx.Value("JWT").(string)
	if !ok || token == "" {
		return nil, errs.BadRequest("missing authorization token")
	}

	res, err := supabase.SupabaseLogout(auth.Client, token)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
