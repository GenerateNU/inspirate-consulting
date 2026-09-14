package routes

import (
	"context"
	"net/http"

	storage "github.com/GenerateNU/inspirate-consulting/internal/data"
	"github.com/GenerateNU/inspirate-consulting/internal/handlers/greeting"
	"github.com/GenerateNU/inspirate-consulting/internal/models"
	"github.com/danielgtaylor/huma/v2"
)

func SetUpGreetingRoutes(api huma.API, store *storage.Store) {
	greetingHandler := greeting.NewHandler(store.Greeting)
	// store := greetingStore.GreetingStore.CreateGreeting()
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
		greeting, err := greetingHandler.CreateGreeting(ctx, input)
		if err != nil {
			return nil, err
		}
		return greeting, nil
	})
}
