package resetpassword

import (
	"context"
	"errors"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"

	"github.com/google/uuid"
)

func (h *Handler) ResetPassword(ctx context.Context, input *models.ResetPasswordInput, supabase supabase.SupabaseInterface) (*models.ResetPasswordResponse, error) {

	supabaseID, ok := ctx.Value("Supabase-ID").(string)
	if !ok {
		return nil, errors.New("could not parse supabase id correctly")
	}

	res, err := supabase.SupabaseResetPassword(auth.Client, input.Body.NewPassword, supabaseID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(supabaseID)
	if err != nil {
		return nil, errs.BadRequest("invalid user ID")
	}

	if err := h.ResetPasswordRepository.UpdateResetTime(ctx, userID, nil); err != nil {
		return nil, err
	}

	return &res, nil
}
