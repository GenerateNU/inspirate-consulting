package greetingRepository

import (
	"context"
	"fmt"

	"inspirate-consulting/internal/models"
)

func (r *GreetingRepository) CreateGreeting(ctx context.Context, greeting models.CreateGreetingInput) (*models.CreateGreetingOutput, error) {
	fmt.Println("You need to create the greeting")
	var fakeCreatedGreeting models.CreateGreetingOutput = models.CreateGreetingOutput{
		Body: models.MessageBody{
			Message: "Hello",
		},
	}

	return &fakeCreatedGreeting, nil

}

