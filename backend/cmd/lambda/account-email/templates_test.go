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
		Name:         "Ada",
		TempPassword: "Xk9mPq2Lz",
		EmailType:    models.AccountCreated,
	}
	subject, body := formatEmail(n, "https://app.example.com/login")

	assert.Contains(t, subject, "Welcome")
	assert.Contains(t, body, "Ada")
	assert.Contains(t, body, "Xk9mPq2Lz")
	assert.Contains(t, body, "https://app.example.com/login")
}

func TestFormatEmail_PasswordReset(t *testing.T) {
	t.Parallel()

	n := models.AccountEmailNotificationRequestBody{
		Name:         "Grace",
		TempPassword: "Vb7qRt3Wy",
		EmailType:    models.PasswordReset,
	}
	subject, body := formatEmail(n, "https://app.example.com/login")

	assert.Contains(t, subject, "reset")
	assert.Contains(t, body, "Grace")
	assert.Contains(t, body, "Vb7qRt3Wy")
	assert.Contains(t, body, "https://app.example.com/login")
}

func TestAccountCreatedEmail_EscapesName(t *testing.T) {
	t.Parallel()

	_, body := accountCreatedEmail(`<script>alert(1)</script>`, "Xk9mPq2Lz", "https://login")

	assert.False(t, strings.Contains(body, "<script>alert(1)</script>"))
	assert.Contains(t, body, "&lt;script&gt;")
}

func TestAccountCreatedEmail_EscapesPassword(t *testing.T) {
	t.Parallel()

	_, body := accountCreatedEmail("Ada", `Xk9<script>&"'`, "https://login")

	assert.False(t, strings.Contains(body, `Xk9<script>&"'`))
	assert.Contains(t, body, "&lt;script&gt;")
}
