package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"log/slog"
	"net/http"
)

func (s *Supabase) SupabaseResetPassword(Client *http.Client, newPassword string, userID string) (models.ResetPasswordResponse, error) {
	payload := map[string]string{"password": newPassword}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return models.ResetPasswordResponse{}, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/auth/v1/admin/users/%s", s.URL, userID), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return models.ResetPasswordResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.ServiceRoleKey))
	req.Header.Set("apikey", s.ServiceRoleKey)

	res, err := Client.Do(req)
	if err != nil {
		slog.Error("Failed to execute Request", "err", err)
		return models.ResetPasswordResponse{}, errs.BadRequest("Failed to execute Request")
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		supabaseError := &models.SupabaseError{}
		slog.Error("Error Response: ", "res.StatusCode", res.StatusCode)
		return models.ResetPasswordResponse{}, errs.NewHTTPError(res.StatusCode, supabaseError)
	}

	return models.ResetPasswordResponse{}, nil
}
