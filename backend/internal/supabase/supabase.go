package supabase

import (
	"net/http"

	"inspirate-consulting/internal/models"
)

type SupabaseInterface interface {
	Signup(email string, password string, Client *http.Client, role string) (models.SignupResponse, error)
	SupabaseLogin(email string, password string, Client *http.Client) (models.LoginResponse, error)
	SupabaseLogout(Client *http.Client, access_token string) (models.LogoutResponse, error)
	SupabaseResetPassword(Client *http.Client, newPassword string, userID string) (models.ResetPasswordResponse, error)
	SupabaseValidateSession(Client *http.Client, access_token string) error
}

// Supabase holds Supabase-related configuration
type Supabase struct {
	URL            string `env:"SUPABASE_URL, required"`
	AnonKey        string `env:"SUPABASE_ANON_KEY, required"`
	ServiceRoleKey string `env:"SUPABASE_SERVICE_ROLE_KEY, required"`
}
