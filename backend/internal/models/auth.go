package models

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

// SignUpPayload represents the payload for Supabase signup
type SignUpPayload struct {
	Email        string            `json:"email" db:"email"`
	Password     string            `json:"password" db:"password"`
	AppMetadata  map[string]string `json:"app_metadata"`
	EmailConfirm bool              `json:"email_confirm"`
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
	ResetTime    *time.Time   `json:"reset_time"`
}

type LogoutInput struct{}

// removed Body to prevent access token and refresh token leak
type LoginOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	ResetTime *time.Time
}

type LogoutResponse struct {
}

type ResetPasswordInput struct {
	Body struct {
		NewPassword string `json:"new_password"`
	}
}

type ResetPasswordResponse struct {
}

// supabaseID here is the student of choice's supabase id
type ForcePasswordResetInput struct {
	SupabaseID uuid.UUID  `path:"supabase_id" required:"true"`
	ResetTime  *time.Time `json:"reset_time"`
}

type ForcePasswordResetOutput struct{}
