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

type ClaimPasswordInput struct {
	Token string `query:"token" required:"true"`
}

type ClaimPasswordBody struct {
	Email        string `json:"email"`
	TempPassword string `json:"temp_password"`
}

type ClaimPasswordOutput struct {
	Body *ClaimPasswordBody
}
