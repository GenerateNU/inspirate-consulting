package user

import (
	"context"
	"inspirate-consulting/internal/models"
)

func (h *Handler) FetchUserBySupabaseID(ctx context.Context, input *models.FetchUserBySupabaseIDInput) (*models.FetchUserOutput, error) {
	user, err := h.UserRepository.FetchUserBySupabaseID(ctx, *input)
	if err != nil {
		return nil, err
	}

	return user, nil
}
