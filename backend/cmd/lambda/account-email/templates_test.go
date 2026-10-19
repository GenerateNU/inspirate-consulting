package main

import (
	"strings"
	"testing"

	"inspirate-consulting/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestFormatEmail_AccountCreated(t *testing.T) {
	t.Parallel()

	n := models.AccountEmailNotificationRequestBody{
		Name:      "Ada",
		EmailType: models.AccountCreated,
	}
	subject, body := formatEmail(n, "https://app.example.com/claim-password?token=abc", "https://app.example.com/login")

	assert.Contains(t, subject, "Welcome")
	assert.Contains(t, body, "Ada")
	assert.Contains(t, body, "https://app.example.com/claim-password?token=abc")
	assert.Contains(t, body, "https://app.example.com/login")
	assert.NotContains(t, body, "temp-pass")
}

func TestFormatEmail_PasswordReset(t *testing.T) {
	t.Parallel()

	n := models.AccountEmailNotificationRequestBody{
		Name:      "Grace",
		EmailType: models.PasswordReset,
	}
	subject, body := formatEmail(n, "https://app.example.com/claim-password?token=xyz", "https://app.example.com/login")

	assert.Contains(t, subject, "reset")
	assert.Contains(t, body, "Grace")
	assert.Contains(t, body, "https://app.example.com/claim-password?token=xyz")
}

func TestAccountCreatedEmail_EscapesName(t *testing.T) {
	t.Parallel()

	_, body := accountCreatedEmail(`<script>alert(1)</script>`, "https://claim", "https://login")

	assert.False(t, strings.Contains(body, "<script>alert(1)</script>"))
	assert.Contains(t, body, "&lt;script&gt;")
}
