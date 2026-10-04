package resetpassword

import (
	"context"
	"time"

	"inspirate-consulting/internal/models"
)

func (h *Handler) ForcePasswordReset(ctx context.Context, input *models.ForcePasswordResetInput) (*models.ForcePasswordResetOutput, error) {
	oneSecondAgo := time.Now().Add(-1 * time.Second)
	if err := h.ResetPasswordRepository.UpdateResetTime(ctx, input.ID, &oneSecondAgo); err != nil {
		return nil, err
	}

	return &models.ForcePasswordResetOutput{}, nil
}
