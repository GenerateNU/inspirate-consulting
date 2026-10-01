package auth

import (
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/models"
)

func SupabaseSignup(s *config.Supabase, email, password string) (models.SignupResponse, error) {
	return s.Signup(email, password, Client)
}
