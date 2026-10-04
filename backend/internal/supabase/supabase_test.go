package supabase

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

func newTestSupabase(url string) *Supabase {
	return &Supabase{
		URL:            url,
		AnonKey:        "test-anon-key",
		ServiceRoleKey: "test-service-role-key",
	}
}

func TestSupabaseLogin_Success(t *testing.T) {
	userID := uuid.New()
	expected := models.LoginResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		User:         models.UserResponse{ID: userID},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	s := newTestSupabase(server.URL)
	resp, err := s.SupabaseLogin("test@example.com", "password123", server.Client())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.AccessToken != expected.AccessToken {
		t.Errorf("expected access token %q, got %q", expected.AccessToken, resp.AccessToken)
	}
	if resp.User.ID != userID {
		t.Errorf("expected user ID %v, got %v", userID, resp.User.ID)
	}
}

func TestSupabaseLogin_InvalidCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.SupabaseError{
			Code:    400,
			Message: "Invalid login credentials",
		})
	}))
	defer server.Close()

	s := newTestSupabase(server.URL)
	_, err := s.SupabaseLogin("bad@example.com", "wrongpassword", server.Client())
	if err == nil {
		t.Fatal("expected an error for invalid credentials, got nil")
	}
}

func TestSupabaseLogout_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	s := newTestSupabase(server.URL)
	_, err := s.SupabaseLogout(server.Client(), "test-access-token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSupabaseLogout_InvalidToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	s := newTestSupabase(server.URL)
	_, err := s.SupabaseLogout(server.Client(), "invalid-token")
	if err == nil {
		t.Fatal("expected an error for invalid token, got nil")
	}
}
