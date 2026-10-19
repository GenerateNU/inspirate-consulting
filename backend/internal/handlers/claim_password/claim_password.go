package claimpassword

import (
	"context"
	"errors"

	"inspirate-consulting/internal/claimtoken"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) ClaimPassword(ctx context.Context, input *models.ClaimPasswordInput) (*models.ClaimPasswordOutput, error) {
	claims, err := claimtoken.Verify(h.Secret, input.Token)
	if err != nil {
		if errors.Is(err, claimtoken.ErrInvalidToken) {
			return nil, errs.Unauthorized()
		}
		return nil, err
	}

	return &models.ClaimPasswordOutput{
		Body: &models.ClaimPasswordBody{
			Email:        claims.Email,
			TempPassword: claims.TempPassword,
		},
	}, nil
}
