package models

import (
	"github.com/google/uuid"
)

// SignUpPayload represents the payload for Supabase signup
type SignUpPayload struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

// UserSignupResponse represents the user data returned from Supabase signup
type UserSignupResponse struct {
	ID uuid.UUID `json:"id"`
}

// SignupResponse represents the complete response from Supabase signup
type SignupResponse struct {
	AccessToken string             `json:"access_token"`
	User        UserSignupResponse `json:"user"`
}

// Error response from Supabase API
type SupabaseError struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"msg"`
}

// Supabase error formatter
func (e *SupabaseError) Error() string {
	return e.Message
}

// Login input schema
type LoginInput struct {
	Body struct {
		Email    string `json:"email" db:"email"`
		Password string `json:"password" db:"password"`
	}
}

// Supabase "User" schema for login
type UserResponse struct {
	ID uuid.UUID `json:"id"`
}

// Login response from Supabase API
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
	Error        interface{}  `json:"error"`
}
