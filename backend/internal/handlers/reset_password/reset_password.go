package resetpassword

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"

	"github.com/google/uuid"
)

func (h *Handler) ResetPassword(ctx context.Context, input *models.ResetPasswordInput, supabase supabase.SupabaseInterface) (*models.ResetPasswordResponse, error) {
	res, err := supabase.SupabaseResetPassword(auth.Client, input.Body.NewPassword, input.UserID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		return nil, errs.BadRequest("invalid user ID")
	}

	if err := h.ResetPasswordRepository.UpdateResetTime(ctx, userID, nil); err != nil {
		return nil, err
	}

	return &res, nil
}
