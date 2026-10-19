package claimtoken

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignVerify_RoundTrip(t *testing.T) {
	t.Parallel()

	token, err := Sign("secret", "student@example.com", "temp-pass-123", time.Minute)
	assert.NoError(t, err)

	claims, err := Verify("secret", token)
	assert.NoError(t, err)
	assert.Equal(t, "student@example.com", claims.Email)
	assert.Equal(t, "temp-pass-123", claims.TempPassword)
}

func TestVerify_Expired(t *testing.T) {
	t.Parallel()

	token, err := Sign("secret", "student@example.com", "temp-pass-123", -time.Minute)
	assert.NoError(t, err)

	_, err = Verify("secret", token)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerify_WrongSecret(t *testing.T) {
	t.Parallel()

	token, err := Sign("secret", "student@example.com", "temp-pass-123", time.Minute)
	assert.NoError(t, err)

	_, err = Verify("different-secret", token)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestVerify_Malformed(t *testing.T) {
	t.Parallel()

	_, err := Verify("secret", "not-a-token")
	assert.ErrorIs(t, err, ErrInvalidToken)
}
