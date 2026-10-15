package routes

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/data/aws/sqsclient"
	"inspirate-consulting/internal/supabase"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

// The guard tests run the real app with TestMode off, so the production
// AuthMiddleware and VerifyRole handle every request. auth.NewVerifier loads
// signing keys from $SUPABASE_URL/auth/v1/.well-known/jwks.json, so we serve a
// test key there and sign tokens with it.

const testKeyID = "test-key"

var (
	testKeyOnce sync.Once
	testKey     *ecdsa.PrivateKey
	tokenSeq    atomic.Int64
)

// signingKey returns the process-wide test key, starting the JWKS server and
// pointing SUPABASE_URL at it on first use (once, so parallel tests agree).
func signingKey() *ecdsa.PrivateKey {
	testKeyOnce.Do(func() {
		testKey = newTestKey()
		jwks := jwksFor(&testKey.PublicKey, testKeyID)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/auth/v1/.well-known/jwks.json" {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jwks)
		}))
		if err := os.Setenv("SUPABASE_URL", srv.URL); err != nil {
			panic(err)
		}
	})
	return testKey
}

func newTestKey() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return key
}

func jwksFor(pub *ecdsa.PublicKey, kid string) []byte {
	b64 := func(n *big.Int) string {
		buf := make([]byte, 32)
		n.FillBytes(buf)
		return base64.RawURLEncoding.EncodeToString(buf)
	}
	out, err := json.Marshal(map[string]any{"keys": []map[string]string{{
		"kty": "EC", "crv": "P-256", "x": b64(pub.X), "y": b64(pub.Y),
		"kid": kid, "alg": "ES256", "use": "sig",
	}}})
	if err != nil {
		panic(err)
	}
	return out
}

// claimsFor returns Supabase-shaped claims for a user with the given
// app_metadata role (no role key when role is "").
func claimsFor(sub, role string) jwt.MapClaims {
	app := map[string]any{"provider": "email"}
	if role != "" {
		app["role"] = role
	}
	now := time.Now()
	return jwt.MapClaims{
		"sub":          sub,
		"aud":          "authenticated",
		"role":         "authenticated",
		"app_metadata": app,
		"iat":          now.Unix(),
		"exp":          now.Add(time.Hour).Unix(),
		// distinct per token so revoking one doesn't affect another
		"session_id": tokenSeq.Add(1),
	}
}

func signWith(method jwt.SigningMethod, key any, kid string, claims jwt.MapClaims) string {
	tok := jwt.NewWithClaims(method, claims)
	tok.Header["kid"] = kid
	signed, err := tok.SignedString(key)
	if err != nil {
		panic(err)
	}
	return signed
}

// sign returns a token the verifier accepts (subject to the claims).
func sign(claims jwt.MapClaims) string {
	return signWith(jwt.SigningMethodES256, signingKey(), testKeyID, claims)
}

// tokenFor returns a valid session token for the given role.
func tokenFor(role string) string {
	return sign(claimsFor("00000000-0000-0000-0000-00000000000"+roleDigit(role), role))
}

func roleDigit(role string) string {
	switch role {
	case roleStudent:
		return "1"
	case roleCounselor:
		return "2"
	default:
		return "3"
	}
}

// authApp is the real app with auth enabled. Every signed token counts as a
// live Supabase session unless it was revoked.
type authApp struct {
	t       *testing.T
	app     *fiber.App
	mu      sync.Mutex
	revoked map[string]bool
}

func newAuthApp(t *testing.T, repo *data.Repository) *authApp {
	t.Helper()
	signingKey()
	a := &authApp{t: t, revoked: map[string]bool{}}
	sb := &supabase.MockSupabase{
		ValidateSessionFunc: func(tok string) error {
			a.mu.Lock()
			defer a.mu.Unlock()
			if a.revoked[tok] {
				return errors.New("session validation failed: status 403")
			}
			return nil
		},
	}
	if repo == nil {
		repo = &data.Repository{}
	}
	app, _, err := SetupApp(config.Config{TestMode: false, Supabase: sb}, repo, &sqsclient.SQSClient{})
	require.NoError(t, err)
	a.app = app
	return a
}

// revoke ends a session on the Supabase side, as sign-out or an admin would.
func (a *authApp) revoke(tok string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.revoked[tok] = true
}

type apiResponse struct {
	Status int
	Body   []byte
}

// Detail is the huma error message ("detail") of an error response.
func (r apiResponse) Detail() string {
	var e struct {
		Detail string `json:"detail"`
	}
	_ = json.Unmarshal(r.Body, &e)
	return e.Detail
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
	return apiResponse{Status: resp.StatusCode, Body: b}
}
