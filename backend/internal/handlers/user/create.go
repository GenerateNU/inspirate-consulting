package user

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
	"log/slog"
)

func (h *Handler) CreateUser(ctx context.Context, input *models.CreateUserInput, supabase supabase.SupabaseInterface) (*models.CreateUserOutput, error) {
	tempPassword, err := generateTempPassword()
	if err != nil {
		return nil, err
	}

	signup_response, err := supabase.Signup(input.Body.Email, tempPassword, auth.Client, string(input.Body.Role))
	if err != nil {
		return nil, err
	}

	user, err := h.UserRepository.CreateUser(ctx, *input, signup_response.User.ID)
	if err != nil {
		return nil, err
	}

	notification := &models.AccountEmailNotificationRequestBody{
		UserEmail:    input.Body.Email,
		Name:         input.Body.Name,
		TempPassword: tempPassword,
		EmailType:    models.AccountCreated,
	}
	if _, err := h.SQSClient.SendAccountEmailNotification(ctx, notification); err != nil {
		slog.Error("failed to queue account-created email", "err", err, "email", input.Body.Email)
	}

	return user, nil
}
