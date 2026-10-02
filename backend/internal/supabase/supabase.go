package supabase

import (
	"net/http"

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
