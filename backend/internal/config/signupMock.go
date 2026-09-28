package config

import (
	"inspirate-consulting/internal/models"
	"net/http"

	"github.com/google/uuid"
)

// mock for supabase
type MockSupabase struct{}

func (s *MockSupabase) Signup(email string, password string, Client *http.Client) (models.SignupResponse, error) {
	user_resp := &models.UserSignupResponse{ID: uuid.New()}

	signup_resp := &models.SignupResponse{AccessToken: "", User: *user_resp}
	return *signup_resp, nil
}
