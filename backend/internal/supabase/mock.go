package supabase

import (
	"inspirate-consulting/internal/models"
	"net/http"

	"github.com/google/uuid"
)

type MockSupabase struct{}

func (s *MockSupabase) Signup(email string, password string, Client *http.Client) (models.SignupResponse, error) {
	return models.SignupResponse{
		AccessToken: "",
		User:        models.UserSignupResponse{ID: uuid.New()},
	}, nil
}

func (s *MockSupabase) SupabaseLogin(email string, password string, Client *http.Client) (models.LoginResponse, error) {
	resp := models.LoginResponse{}
	return resp, nil
}

func (s *MockSupabase) SupabaseLogout(Client *http.Client, access_token string) (models.LogoutResponse, error) {
	resp := models.LogoutResponse{}
	return resp, nil
}

func (s *MockSupabase) SupabaseResetPassword(Client *http.Client, access_token string) (models.ResetPasswordResponse, error) {
	resp := models.ResetPasswordResponse{}
	return resp, nil
}
