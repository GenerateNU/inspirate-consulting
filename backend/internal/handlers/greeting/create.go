package greeting

import (
	"context"
	"fmt"

	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateGreeting(ctx context.Context, input *models.CreateGreetingInput) (*models.CreateGreetingOutput, error) {
	name := input.Name
	fmt.Println("At CreateGreeting handler")
	// Note we would catch this error way earlier normally
	if len(name) < 2 {
		return nil, fmt.Errorf("Name '%s' must be longer than 1 character")
	}
	greetingWithName := fmt.Sprintf("Hello, %s!", name)
	greetingResponse := models.CreateGreetingOutput{
		Body: models.MessageBody{
			Message: greetingWithName,
		},
	}
	return &greetingResponse, nil
}
