package login

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
	"net/http"
)

func (h *Handler) Login(ctx context.Context, input *models.LoginInput, supabase supabase.SupabaseInterface) (*models.LoginOutput, error) {
	res, err := supabase.SupabaseLogin(input.Body.Email, input.Body.Password, auth.Client)
	if err != nil {
		return nil, err
	}

	userInput := models.FetchUserBySupabaseIDInput{}
	userInput.SupabaseID = res.User.ID
	user, err := h.LoginRepository.FetchUserBySupabaseID(ctx, userInput)
	if err != nil {
		return nil, err
	}

	res.ResetTime = user.Body.ResetTime
	return &models.LoginOutput{
		SetCookie: http.Cookie{
			Name:     "jwt",
			Value:    res.AccessToken,
			Path:     "/",
			MaxAge:   res.ExpiresIn,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		},
		Body: &res,
	}, nil
}
