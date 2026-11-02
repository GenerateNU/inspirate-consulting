package models

type AccountEmailType string

const (
	AccountCreated AccountEmailType = "account_created"
	PasswordReset  AccountEmailType = "password_reset"
)

type AccountEmailNotificationRequestBody struct {
	UserEmail    string           `json:"user_email"`
	Name         string           `json:"name"`
	TempPassword string           `json:"temp_password"`
	EmailType    AccountEmailType `json:"email_type"`
}
