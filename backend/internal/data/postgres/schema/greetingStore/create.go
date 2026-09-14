package greetingStore

import (
	"context"
	"fmt"

	"github.com/GenerateNU/inspirate-consulting/internal/models"
)

func (r *GreetingStore) CreateGreeting(ctx context.Context, greeting models.CreateGreetingInput) (*models.CreateGreetingOutput, error) {
	fmt.Println("You need to create the greeting")
	var fakeCreatedGreeting models.CreateGreetingOutput = models.CreateGreetingOutput{
		Body: models.MessageBody{
			Message: "Hello",
		},
	}

	return &fakeCreatedGreeting, nil

}
