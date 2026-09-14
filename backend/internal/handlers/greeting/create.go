package greeting

import (
	"context"
	"fmt"

	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateGreeting(ctx context.Context, input *models.CreateGreetingInput) (*models.CreateGreetingOutput, error) {
	fmt.Println("Made it to the CreateGreeting handler")
	greeting, err := h.GreetingRepository.CreateGreeting(ctx, *input)
	if err != nil {
		return nil, err
	}
	return greeting, nil
}
