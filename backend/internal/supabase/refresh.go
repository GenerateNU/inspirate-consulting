package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"io"
	"log/slog"
	"net/http"
)

func (s *Supabase) SupabaseRefreshToken(Client *http.Client, refreshToken string) (models.LoginResponse, error) {
	payload := models.RefreshTokenPayload{RefreshToken: refreshToken}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return models.LoginResponse{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/token?grant_type=refresh_token", s.URL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return models.LoginResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.ServiceRoleKey))
	req.Header.Set("apikey", s.ServiceRoleKey)

	res, err := Client.Do(req)
	if err != nil {
		slog.Error("Failed to execute Request", "err", err)
		return models.LoginResponse{}, errs.BadRequest("Failed to execute Request")
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Failed to read response body", "err", err)
		return models.LoginResponse{}, errs.BadRequest("Failed to read response body")
	}

	if res.StatusCode != http.StatusOK {
		supabaseError := &models.SupabaseError{}
		if err := json.Unmarshal(body, supabaseError); err != nil {
			slog.Error("Error parsing response: ", "err", err)
			return models.LoginResponse{}, err
		}
		slog.Error("Error Response: ", "res.StatusCode", res.StatusCode, "body", string(body))
		return models.LoginResponse{}, errs.NewHTTPError(res.StatusCode, supabaseError)
	}

	var refreshResponse models.LoginResponse
	if err := json.Unmarshal(body, &refreshResponse); err != nil {
		slog.Error("Failed to parse response body", "body", err)
		return models.LoginResponse{}, errs.BadRequest("Failed to parse response body")
	}

	if refreshResponse.Error != nil {
		return models.LoginResponse{}, errs.BadRequest(fmt.Sprintf("Refresh Token Response Error %v", refreshResponse.Error))
	}

	return refreshResponse, nil
}
