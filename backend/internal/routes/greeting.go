package routes

import (
	"context"
	"fmt"
	"net/http"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/greeting"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpGreetingRoutes(api huma.API, repository *data.Repository) {
	greetingHandler := greeting.NewHandler(repository.Greeting)
	// repository := greetingRepository.GreetingRepository.CreateGreeting()
	// Register GET /greeting/{name} handler.
	// The handler function takes in a struct that defines its inputs ('name' in this case)
	// and returns the CreateGreetingOutput model built in the models
	huma.Register(api, huma.Operation{
		OperationID: "get-greeting",
		Method:      http.MethodGet,
		Path:        "/greeting/{name}",
		Description: "Get a greeting for a person by name.",
		Tags:        []string{"Greetings"},
	}, func(ctx context.Context, input *models.CreateGreetingInput) (*models.CreateGreetingOutput, error) {
		fmt.Println("Made it to router")
		greeting, err := greetingHandler.CreateGreeting(ctx, input)
		if err != nil {
			return nil, err
		}
		return greeting, nil
	})
}
