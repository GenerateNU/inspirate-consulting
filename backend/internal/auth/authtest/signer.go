// Package authtest provides test doubles for exercising the real auth stack
// (AuthMiddleware + VerifyRole) without a running Supabase instance.
//
// It is only meant to be imported from _test.go files.
package authtest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const KeyID = "authtest-key"

// Signer mints Supabase-shaped JWTs signed with a key that is published on a
// local JWKS endpoint, so auth.NewVerifier accepts them like real Supabase tokens.
type Signer struct {
	key    *ecdsa.PrivateKey
	server *httptest.Server
	seq    atomic.Int64
}

var (
	sharedOnce   sync.Once
	sharedSigner *Signer
)

// SharedSigner returns one process-wide signer and points SUPABASE_URL at its
// JWKS endpoint (routes.SetupApp reads that env var when building the verifier).
// A single shared instance keeps parallel tests from fighting over the env var.
func SharedSigner() *Signer {
	sharedOnce.Do(func() {
		sharedSigner = newSigner()
		if err := os.Setenv("SUPABASE_URL", sharedSigner.server.URL); err != nil {
			panic(err)
		}
	})
	return sharedSigner
}

func newSigner() *Signer {
	key := NewKey()
	s := &Signer{key: key}
	jwks := jwksJSON(&key.PublicKey, KeyID)
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/v1/.well-known/jwks.json" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(jwks)
	}))
	return s
}

// NewKey generates a fresh ES256 key, e.g. to sign a token the verifier must reject.
func NewKey() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return key
}

// Claims are the Supabase claims a token is minted with.
type Claims struct {
	Sub          string
	Email        string
	AppMetadata  map[string]any
	UserMetadata map[string]any
	IssuedAt     time.Time
	ExpiresAt    time.Time
	NotBefore    time.Time
}

// ClaimsFor returns valid claims for a user with the given app_metadata role.
// An empty role leaves app_metadata without a role key.
func ClaimsFor(sub, role string) Claims {
	app := map[string]any{"provider": "email"}
	if role != "" {
		app["role"] = role
	}
	now := time.Now()
	return Claims{
		Sub:         sub,
		Email:       sub + "@example.com",
		AppMetadata: app,
		IssuedAt:    now,
		ExpiresAt:   now.Add(time.Hour),
	}
}

func (c Claims) mapClaims(jti int64) jwt.MapClaims {
	m := jwt.MapClaims{
		"sub":           c.Sub,
		"email":         c.Email,
		"aud":           "authenticated",
		"role":          "authenticated",
		"app_metadata":  c.AppMetadata,
		"user_metadata": c.UserMetadata,
		"iat":           c.IssuedAt.Unix(),
		"exp":           c.ExpiresAt.Unix(),
		// unique per token so two logins of the same user get distinct sessions
		"session_id": fmt.Sprintf("session-%d", jti),
	}
	if !c.NotBefore.IsZero() {
		m["nbf"] = c.NotBefore.Unix()
	}
	return m
}

// Sign mints a token the verifier accepts (subject to the claims' validity).
func (s *Signer) Sign(c Claims) string {
	return SignWith(s.key, KeyID, c, s.seq.Add(1))
}

// SignWith signs claims with an arbitrary key/kid, for forged-token tests.
func SignWith(key *ecdsa.PrivateKey, kid string, c Claims, jti int64) string {
	tok := jwt.NewWithClaims(jwt.SigningMethodES256, c.mapClaims(jti))
	tok.Header["kid"] = kid
	signed, err := tok.SignedString(key)
	if err != nil {
		panic(err)
	}
	return signed
}

// SignHS256 signs claims with a shared secret; the verifier only allows ES256/RS256.
func SignHS256(secret []byte, c Claims) string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, c.mapClaims(0))
	tok.Header["kid"] = KeyID
	signed, err := tok.SignedString(secret)
	if err != nil {
		panic(err)
	}
	return signed
}

// SignNone builds an unsigned ("alg": "none") token.
func SignNone(c Claims) string {
	tok := jwt.NewWithClaims(jwt.SigningMethodNone, c.mapClaims(0))
	tok.Header["kid"] = KeyID
	signed, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		panic(err)
	}
	return signed
}

func jwksJSON(pub *ecdsa.PublicKey, kid string) []byte {
	b64 := func(n *big.Int) string {
		buf := make([]byte, 32)
		n.FillBytes(buf)
		return base64.RawURLEncoding.EncodeToString(buf)
	}
	out, err := json.Marshal(map[string]any{
		"keys": []map[string]string{{
			"kty": "EC",
			"crv": "P-256",
			"x":   b64(pub.X),
			"y":   b64(pub.Y),
			"kid": kid,
			"alg": "ES256",
			"use": "sig",
		}},
	})
	if err != nil {
		panic(err)
	}
	return out
}
