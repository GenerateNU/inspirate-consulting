package resetpassword

import (
	"context"
	"log/slog"
	"time"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/passwordgen"
	"inspirate-consulting/internal/supabase"
)

func (h *Handler) SendTempPassword(ctx context.Context, input *models.SendTempPasswordInput, sb supabase.SupabaseInterface) (*models.SendTempPasswordOutput, error) {
	user, err := h.ResetPasswordRepository.FetchUser(ctx, models.FetchUserInput{ID: input.ID})
	if err != nil {
		return nil, err
	}

	tempPassword, err := passwordgen.Generate()
	if err != nil {
		return nil, err
	}

	if _, err := sb.SupabaseResetPassword(auth.Client, tempPassword, user.Body.SupabaseID.String()); err != nil {
		return nil, err
	}

	oneSecondAgo := time.Now().Add(-1 * time.Second)
	if err := h.ResetPasswordRepository.UpdateResetTime(ctx, user.Body.SupabaseID, &oneSecondAgo); err != nil {
		return nil, err
	}

	notification := &models.AccountEmailNotificationRequestBody{
		UserEmail:    user.Body.Email,
		Name:         user.Body.Name,
		TempPassword: tempPassword,
		EmailType:    models.PasswordReset,
	}
	if _, err := h.SQSClient.SendAccountEmailNotification(ctx, notification); err != nil {
		slog.Error("failed to queue password-reset email", "err", err, "user_id", input.ID)
	}

	return &models.SendTempPasswordOutput{
		Body: &models.SendTempPasswordBody{Email: user.Body.Email},
	}, nil
}
