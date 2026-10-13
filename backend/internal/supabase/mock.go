package supabase

import (
	"inspirate-consulting/internal/models"
	"net/http"

	"github.com/google/uuid"
)

// MockSupabase is a test double for SupabaseInterface. Each call succeeds with
// an empty response unless the matching *Func override is set.
type MockSupabase struct {
	SignupFunc          func(email, password, role string) (models.SignupResponse, error)
	LoginFunc           func(email, password string) (models.LoginResponse, error)
	LogoutFunc          func(accessToken string) (models.LogoutResponse, error)
	ResetPasswordFunc   func(newPassword, userID string) (models.ResetPasswordResponse, error)
	ValidateSessionFunc func(accessToken string) error
}

func (s *MockSupabase) Signup(email string, password string, Client *http.Client, role string) (models.SignupResponse, error) {
	if s.SignupFunc != nil {
		return s.SignupFunc(email, password, role)
	}
	return models.SignupResponse{
		AccessToken: "",
		User:        models.UserSignupResponse{ID: uuid.New()},
	}, nil
}

func (s *MockSupabase) SupabaseLogin(email string, password string, Client *http.Client) (models.LoginResponse, error) {
	if s.LoginFunc != nil {
		return s.LoginFunc(email, password)
	}
	resp := models.LoginResponse{}
	return resp, nil
}

func (s *MockSupabase) SupabaseLogout(Client *http.Client, access_token string) (models.LogoutResponse, error) {
	if s.LogoutFunc != nil {
		return s.LogoutFunc(access_token)
	}
	resp := models.LogoutResponse{}
	return resp, nil
}

func (s *MockSupabase) SupabaseResetPassword(Client *http.Client, access_token string, userID string) (models.ResetPasswordResponse, error) {
	if s.ResetPasswordFunc != nil {
		return s.ResetPasswordFunc(access_token, userID)
	}
	resp := models.ResetPasswordResponse{}
	return resp, nil
}

func (s *MockSupabase) SupabaseValidateSession(Client *http.Client, access_token string) error {
	if s.ValidateSessionFunc != nil {
		return s.ValidateSessionFunc(access_token)
	}
	return nil
}
