package authtest

import (
	"errors"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

var _ supabase.SupabaseInterface = (*FakeSupabase)(nil)

// Account is a Supabase auth user held by FakeSupabase.
type Account struct {
	ID       uuid.UUID
	Email    string
	Password string
	Role     string
}

// FakeSupabase is an in-memory stand-in for Supabase Auth. Login mints real
// signed JWTs (via Signer) carrying the account's role in app_metadata, and
// sessions can be revoked, so a request travels the same path it would in prod.
//
// The *Err fields force the matching call to fail, for error-path tests.
type FakeSupabase struct {
	Signer *Signer

	SignupErr        error
	LoginErr         error
	LogoutErr        error
	ResetPasswordErr error

	mu       sync.Mutex
	accounts map[string]*Account // by email
	sessions map[string]bool     // token -> live

	SignupCalls        []SignupCall
	LogoutCalls        []string
	ResetPasswordCalls []ResetPasswordCall
}

type SignupCall struct {
	Email, Password, Role string
}

type ResetPasswordCall struct {
	UserID, NewPassword string
}

func NewFakeSupabase(signer *Signer) *FakeSupabase {
	return &FakeSupabase{
		Signer:   signer,
		accounts: map[string]*Account{},
		sessions: map[string]bool{},
	}
}

// AddAccount registers an auth user directly (as if created earlier).
func (f *FakeSupabase) AddAccount(email, password, role string) *Account {
	f.mu.Lock()
	defer f.mu.Unlock()
	a := &Account{ID: uuid.New(), Email: email, Password: password, Role: role}
	f.accounts[email] = a
	return a
}

// Account returns the auth user for email, if any.
func (f *FakeSupabase) Account(email string) (*Account, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	a, ok := f.accounts[email]
	return a, ok
}

// IssueSession mints a live session token for a subject/role without a login,
// for tests that only care about what a given role can reach.
func (f *FakeSupabase) IssueSession(sub, role string) string {
	return f.IssueSessionWithClaims(ClaimsFor(sub, role))
}

// IssueSessionWithClaims registers a live session for a token with custom claims.
func (f *FakeSupabase) IssueSessionWithClaims(c Claims) string {
	tok := f.Signer.Sign(c)
	f.mu.Lock()
	f.sessions[tok] = true
	f.mu.Unlock()
	return tok
}

// RegisterSession marks an arbitrary token as live on the Supabase side, so a
// test can isolate the JWT checks from the session check.
func (f *FakeSupabase) RegisterSession(tok string) {
	f.mu.Lock()
	f.sessions[tok] = true
	f.mu.Unlock()
}

// Revoke ends a session, as Supabase does on sign-out or admin revocation.
func (f *FakeSupabase) Revoke(tok string) {
	f.mu.Lock()
	delete(f.sessions, tok)
	f.mu.Unlock()
}

func (f *FakeSupabase) Signup(email string, password string, _ *http.Client, role string) (models.SignupResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.SignupCalls = append(f.SignupCalls, SignupCall{Email: email, Password: password, Role: role})
	if f.SignupErr != nil {
		return models.SignupResponse{}, f.SignupErr
	}
	if _, exists := f.accounts[email]; exists {
		return models.SignupResponse{}, errs.NewHTTPError(http.StatusUnprocessableEntity,
			errors.New("A user with this email address has already been registered"))
	}
	a := &Account{ID: uuid.New(), Email: email, Password: password, Role: role}
	f.accounts[email] = a
	return models.SignupResponse{User: models.UserSignupResponse{ID: a.ID}}, nil
}

func (f *FakeSupabase) SupabaseLogin(email string, password string, _ *http.Client) (models.LoginResponse, error) {
	if f.LoginErr != nil {
		return models.LoginResponse{}, f.LoginErr
	}
	f.mu.Lock()
	a, ok := f.accounts[email]
	f.mu.Unlock()
	if !ok || a.Password != password {
		return models.LoginResponse{}, errs.NewHTTPError(http.StatusBadRequest, errors.New("Invalid login credentials"))
	}
	c := ClaimsFor(a.ID.String(), a.Role)
	c.Email = a.Email
	tok := f.IssueSessionWithClaims(c)
	return models.LoginResponse{
		AccessToken:  tok,
		TokenType:    "bearer",
		ExpiresIn:    3600,
		RefreshToken: "refresh-" + tok[len(tok)-8:],
		User:         models.UserResponse{ID: a.ID},
	}, nil
}

func (f *FakeSupabase) SupabaseLogout(_ *http.Client, accessToken string) (models.LogoutResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.LogoutCalls = append(f.LogoutCalls, accessToken)
	if f.LogoutErr != nil {
		return models.LogoutResponse{}, f.LogoutErr
	}
	delete(f.sessions, accessToken)
	return models.LogoutResponse{}, nil
}

func (f *FakeSupabase) SupabaseResetPassword(_ *http.Client, newPassword string, userID string) (models.ResetPasswordResponse, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ResetPasswordCalls = append(f.ResetPasswordCalls, ResetPasswordCall{UserID: userID, NewPassword: newPassword})
	if f.ResetPasswordErr != nil {
		return models.ResetPasswordResponse{}, f.ResetPasswordErr
	}
	for _, a := range f.accounts {
		if a.ID.String() == userID {
			a.Password = newPassword
			return models.ResetPasswordResponse{}, nil
		}
	}
	return models.ResetPasswordResponse{}, errs.NewHTTPError(http.StatusNotFound, errors.New("User not found"))
}

func (f *FakeSupabase) SupabaseValidateSession(_ *http.Client, accessToken string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.sessions[accessToken] {
		return errors.New("session validation failed: status 403")
	}
	return nil
}
