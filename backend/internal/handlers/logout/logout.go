package logout

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
	"net/http"
)

func (h *Handler) Logout(ctx context.Context, input *models.LogoutInput, supabase supabase.SupabaseInterface) (*models.LogoutResponse, error) {
	token, ok := ctx.Value("JWT").(string)
	if !ok || token == "" {
		return nil, errs.BadRequest("missing authorization token")
	}

	if _, err := supabase.SupabaseLogout(auth.Client, token); err != nil {
		return nil, err
	}

	// resets the cookies
	return &models.LogoutResponse{
		SetCookie: []http.Cookie{
			auth.ExpiredAccessTokenCookie(),
			auth.ExpiredRefreshTokenCookie(),
		},
	}, nil
}
