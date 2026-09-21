package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/greeting"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func CreateUser(api huma.API, repository *data.Repository) {
	greetingHandler := greeting.NewHandler(repository.Greeting)
	// Register POST /greeting handler.
	// The handler function takes in a struct that defines its inputs ('Body' in this case)
	// and returns the CreateGreetingOutput model built in the models
	huma.Register(api, huma.Operation{
		OperationID: "create-greeting",
		Method:      http.MethodPost,
		Path:        "/greeting",
		Description: "Create a greeting for a person by name.",
		Tags:        []string{"Greetings"},
	}, func(ctx context.Context, input *models.CreateGreetingInput) (*models.CreateGreetingOutput, error) {
		greeting, err := greetingHandler.CreateGreeting(ctx, input)
		if err != nil {
			return nil, err
		}
		return greeting, nil
	})
}
