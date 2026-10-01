package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

type SupabaseInterface interface {
	Signup(email string, password string, Client *http.Client) (models.SignupResponse, error)
}

// Supabase holds Supabase-related configuration
type Supabase struct {
	URL            string `env:"SUPABASE_URL, required"`
	AnonKey        string `env:"SUPABASE_ANON_KEY, required"`
	ServiceRoleKey string `env:"SUPABASE_SERVICE_ROLE_KEY, required"`
}

func (s *Supabase) Signup(email string, password string, Client *http.Client) (models.SignupResponse, error) {
	if err := validatePasswordStrength(password); err != nil {
		return models.SignupResponse{}, err
	}

	payload := models.SignUpPayload{Email: email, Password: password}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return models.SignupResponse{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/signup", s.URL), bytes.NewBuffer(payloadBytes))
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

	var response models.SignupResponse
	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("Error parsing response: ", "err", err)
		return models.SignupResponse{}, err
	}

	return response, nil
}

func (s *Supabase) SupabaseLogin(email string, password string, Client *http.Client) (models.LoginResponse, error) {
	payload := models.SignUpPayload{Email: email, Password: password}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return models.LoginResponse{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/token?grant_type=password", s.URL), bytes.NewBuffer(payloadBytes))
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

	var signInResponse models.LoginResponse
	if err := json.Unmarshal(body, &signInResponse); err != nil {
		slog.Error("Failed to parse response body", "body", err)
		return models.LoginResponse{}, errs.BadRequest("Failed to parse response body")
	}

	if signInResponse.Error != nil {
		return models.LoginResponse{}, errs.BadRequest(fmt.Sprintf("Sign In Response Error %v", signInResponse.Error))
	}

	return signInResponse, nil
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
