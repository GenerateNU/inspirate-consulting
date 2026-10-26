package claimpassword

import (
	"context"
	"testing"
	"time"

	"inspirate-consulting/internal/claimtoken"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestClaimPassword_Success(t *testing.T) {
	t.Parallel()

	token, err := claimtoken.Sign("secret", "student@example.com", "temp-pass-123", time.Minute)
	assert.NoError(t, err)

	handler := NewHandler("secret")
	res, err := handler.ClaimPassword(context.Background(), &models.ClaimPasswordInput{Token: token})

	assert.NoError(t, err)
	assert.Equal(t, "student@example.com", res.Body.Email)
	assert.Equal(t, "temp-pass-123", res.Body.TempPassword)
}

func TestClaimPassword_InvalidToken(t *testing.T) {
	t.Parallel()

	handler := NewHandler("secret")
	res, err := handler.ClaimPassword(context.Background(), &models.ClaimPasswordInput{Token: "garbage"})

	assert.Nil(t, res)
	assert.Equal(t, errs.Unauthorized(), err)
}

func TestClaimPassword_ExpiredToken(t *testing.T) {
	t.Parallel()

	token, err := claimtoken.Sign("secret", "student@example.com", "temp-pass-123", -time.Minute)
	assert.NoError(t, err)

	handler := NewHandler("secret")
	res, err := handler.ClaimPassword(context.Background(), &models.ClaimPasswordInput{Token: token})

	assert.Nil(t, res)
	assert.Equal(t, errs.Unauthorized(), err)
}
