package routes

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/data/aws/sqsclient"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// The guard tests run the real app with TestMode off, so the production
// AuthMiddleware and VerifyRole handle every request. auth.NewVerifier loads
// signing keys from $SUPABASE_URL/auth/v1/.well-known/jwks.json, so we serve a
// test key there and sign tokens with it.

const testKeyID = "test-key"

var (
	testKeyOnce sync.Once
	testKey     *rsa.PrivateKey
	tokenSeq    atomic.Int64
)

// signingKey returns the process-wide test key, starting the JWKS server and
// pointing SUPABASE_URL at it on first use (once, so parallel tests agree).
func signingKey() *rsa.PrivateKey {
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

func newTestKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return key
}

func jwksFor(pub *rsa.PublicKey, kid string) []byte {
	b64 := func(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }
	out, err := json.Marshal(map[string]any{"keys": []map[string]string{{
		"kty": "RSA", "n": b64(pub.N.Bytes()), "e": b64(big.NewInt(int64(pub.E)).Bytes()),
		"kid": kid, "alg": "RS256", "use": "sig",
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
	return signWith(jwt.SigningMethodRS256, signingKey(), testKeyID, claims)
}

// tokenFor returns a valid session token for the given role.
func tokenFor(role string) string {
	return sign(claimsFor(subFor(role), role))
}

// subFor is the Supabase ID of the seeded user that holds the given role.
func subFor(role string) string {
	return "00000000-0000-0000-0000-00000000000" + roleDigit(role)
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

// authApp is the real app with auth enabled, backed by a fresh test database.
// Every signed token counts as a live Supabase session unless it was revoked.
type authApp struct {
	t       *testing.T
	app     *fiber.App
	mu      sync.Mutex
	revoked map[string]bool
}

// newAuthApp seeds one user per role (see subFor) so ResolveIdentityMiddleware
// finds them; override swaps in mocked repositories where a test needs them.
func newAuthApp(t *testing.T, override func(*data.Repository)) *authApp {
	t.Helper()
	signingKey()
	db := testutils.SetupTestDB(t)
	seedIdentities(t, db)
	repo := data.NewRepository(db)
	if override != nil {
		override(repo)
	}
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
	app, _, err := SetupApp(config.Config{TestMode: false, Supabase: sb}, repo, &sqsclient.SQSClient{})
	require.NoError(t, err)
	a.app = app
	return a
}

// seedIdentities inserts a student, a counselor, and a user with no profile
// (for roles the app doesn't know), each linked to subFor(role).
func seedIdentities(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	insertUser := func(role string) string {
		var id string
		err := db.QueryRow(ctx,
			`INSERT INTO users (name, email, supabase_id) VALUES ($1, $2, $3) RETURNING id`,
			role, role+"@example.com", subFor(role)).Scan(&id)
		require.NoError(t, err)
		return id
	}

	var counselorID string
	err := db.QueryRow(ctx, `INSERT INTO counselor (user_id) VALUES ($1) RETURNING id`,
		insertUser(roleCounselor)).Scan(&counselorID)
	require.NoError(t, err)

	_, err = db.Exec(ctx,
		`INSERT INTO student (user_id, year, gpa, review_balance, counselor_id) VALUES ($1, 'senior', 4, 0, $2)`,
		insertUser(roleStudent), counselorID)
	require.NoError(t, err)

	insertUser("admin")
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
