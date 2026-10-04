package supabase

import (
	"inspirate-consulting/internal/models"
	"net/http"
)

func (s *Supabase) SupabaseResetPassword(Client *http.Client, access_token string) (models.ResetPasswordResponse, error) {
	return models.ResetPasswordResponse{}, nil

}
