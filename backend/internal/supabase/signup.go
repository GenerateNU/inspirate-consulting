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
	"regexp"
)

func (s *Supabase) Signup(email string, password string, Client *http.Client, role string) (models.SignupResponse, error) {
	if err := validatePasswordStrength(password); err != nil {
		return models.SignupResponse{}, err
	}

	app_metadata := make(map[string]string)
	app_metadata["role"] = role

	payload := models.SignUpPayload{Email: email, Password: password, AppMetadata: app_metadata}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return models.SignupResponse{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/admin/users", s.URL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		slog.Error("Error in Request Creation: ", "err", err)
		return models.SignupResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.ServiceRoleKey))
	req.Header.Set("apikey", s.ServiceRoleKey)

	res, err := Client.Do(req)
	if err != nil {
		slog.Error("Error executing request: ", "err", err)
		return models.SignupResponse{}, err
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error reading response body: ", "body", body)
		return models.SignupResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		supabaseError := &models.SupabaseError{}
		if err := json.Unmarshal(body, supabaseError); err != nil {
			slog.Error("Error parsing response: ", "err", err)
			return models.SignupResponse{}, err
		}
		slog.Error("Error Response: ", "res.StatusCode", res.StatusCode, "body", string(body))
		return models.SignupResponse{}, errs.NewHTTPError(res.StatusCode, supabaseError)
	}

	// The admin endpoint returns the user object at the top level, not nested under "user"
	var response models.SignupResponse
	if err := json.Unmarshal(body, &response.User); err != nil {
		slog.Error("Error parsing response: ", "err", err)
		return models.SignupResponse{}, err
	}

	return response, nil
}

func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errs.BadRequest("Password must be at least 8 characters long")
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#~$%^&*()+|_.,;<>?/{}\-]`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errs.BadRequest("Password must include uppercase, lowercase, digit and special characters")
	}

	return nil
}
