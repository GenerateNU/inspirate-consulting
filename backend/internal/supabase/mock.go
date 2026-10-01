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
