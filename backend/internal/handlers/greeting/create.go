package greeting

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateGreeting(ctx context.Context, input *models.CreateGreetingInput) (*models.CreateGreetingOutput, error) {
	greeting, err := h.GreetingRepository.CreateGreeting(ctx, *input)
	if err != nil {
		return nil, err
	}
	return greeting, nil
}
