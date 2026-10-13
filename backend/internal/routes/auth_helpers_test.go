package routes

import (
	"bytes"
	"encoding/json"
	"inspirate-consulting/internal/auth/authtest"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

// authApp is the real app (TestMode off) wired to fake Supabase auth, so the
// production AuthMiddleware and VerifyRole guards run on every request.
type authApp struct {
	t     *testing.T
	app   *fiber.App
	sb    *authtest.FakeSupabase
	users *authtest.MemUserRepo
}

// newAuthApp builds the real app with auth enabled. repo may be nil; its User
// repository is always the in-memory one so login/create/reset work end to end.
func newAuthApp(t *testing.T, repo *data.Repository) *authApp {
	t.Helper()
	signer := authtest.SharedSigner() // also points SUPABASE_URL at the fake JWKS
	sb := authtest.NewFakeSupabase(signer)
	users := authtest.NewMemUserRepo()
	if repo == nil {
		repo = &data.Repository{}
	}
	repo.User = users

	app, _, err := SetupApp(config.Config{TestMode: false, Supabase: sb}, repo)
	require.NoError(t, err)
	return &authApp{t: t, app: app, sb: sb, users: users}
}

type apiResponse struct {
	Status  int
	Body    []byte
	Cookies []*http.Cookie
	Header  http.Header
}

func (r apiResponse) JSON(t *testing.T, v any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(r.Body, v), "body: %s", r.Body)
}

// Detail is the huma error message ("detail") of an error response.
func (r apiResponse) Detail() string {
	var e struct {
		Detail string `json:"detail"`
	}
	_ = json.Unmarshal(r.Body, &e)
	return e.Detail
}

func (r apiResponse) Cookie(name string) *http.Cookie {
	for _, c := range r.Cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// do sends a request with the jwt cookie set to token (no cookie if token is "").
func (a *authApp) do(method, path, token string, body any) apiResponse {
	a.t.Helper()
	return a.doReq(method, path, body, func(req *http.Request) {
		if token != "" {
			req.AddCookie(&http.Cookie{Name: "jwt", Value: token})
		}
	})
}

func (a *authApp) doReq(method, path string, body any, mutate func(*http.Request)) apiResponse {
	a.t.Helper()
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(a.t, err)
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, path, reader)
	require.NoError(a.t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if mutate != nil {
		mutate(req)
	}
	resp, err := a.app.Test(req, fiber.TestConfig{Timeout: 10 * time.Second})
	require.NoError(a.t, err)
	defer func() { _ = resp.Body.Close() }()
	b, err := io.ReadAll(resp.Body)
	require.NoError(a.t, err)
	return apiResponse{Status: resp.StatusCode, Body: b, Cookies: resp.Cookies(), Header: resp.Header}
}

// login posts credentials and returns the response plus the session token from
// the jwt cookie ("" when login failed).
func (a *authApp) login(email, password string) (apiResponse, string) {
	a.t.Helper()
	resp := a.do(http.MethodPost, "/user/login", "", map[string]string{"email": email, "password": password})
	if c := resp.Cookie("jwt"); c != nil {
		return resp, c.Value
	}
	return resp, ""
}

// mustLogin logs in and fails the test unless it succeeded.
func (a *authApp) mustLogin(email, password string) string {
	a.t.Helper()
	resp, tok := a.login(email, password)
	require.Equal(a.t, http.StatusOK, resp.Status, "login %s: %s", email, resp.Body)
	require.NotEmpty(a.t, tok, "login %s set no jwt cookie", email)
	return tok
}

// seedUser creates a Supabase account plus its users row (no pending reset).
func (a *authApp) seedUser(name, email, password, role string) *authtest.Account {
	a.t.Helper()
	acct := a.sb.AddAccount(email, password, role)
	a.users.Seed(name, acct.ID)
	return acct
}
