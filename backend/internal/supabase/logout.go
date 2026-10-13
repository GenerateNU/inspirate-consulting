package supabase

import (
	"fmt"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"log/slog"
	"net/http"
)

func (s *Supabase) SupabaseLogout(Client *http.Client, access_token string) (models.LogoutResponse, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/logout", s.URL), nil)
	if err != nil {
		return models.LogoutResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	req.Header.Set("apikey", s.AnonKey)

	res, err := Client.Do(req)
	if err != nil {
		slog.Error("Failed to execute Request", "err", err)
		return models.LogoutResponse{}, errs.BadRequest("Failed to execute Request")
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusNoContent {
		supabaseError := &models.SupabaseError{}
		slog.Error("Error Response: ", "res.StatusCode", res.StatusCode)
		return models.LogoutResponse{}, errs.NewHTTPError(res.StatusCode, supabaseError)
	}

	var logoutResponse models.LogoutResponse

	return logoutResponse, nil

}
