package supabase

import (
	"inspirate-consulting/internal/models"
	"net/http"

	"github.com/google/uuid"
)

type MockClient struct{}

func (s *MockClient) Signup(email string, password string, httpClient *http.Client) (models.SignupResponse, error) {
	return models.SignupResponse{
		AccessToken: "",
		User:        models.UserSignupResponse{ID: uuid.New()},
	}, nil
}
