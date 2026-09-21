package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/user"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetupUserRoutes(api huma.API, repository *data.Repository, config *config.Config) {
	userHandler := user.NewHandler(repository.User)
	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/user",
		Description: "Create a user (student/counselor)",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *models.CreateUserInput) (*models.CreateUserOutput, error) {
		userOutput, err := userHandler.CreateUser(ctx, input, config)
<<<<<<< HEAD
		if err != nil {
			return nil, err
		}
		return userOutput, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "fetch-user",
		Method:      http.MethodGet,
		Path:        "/user/{id}",
		Description: "fetch a user (student/counselor)",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *models.FetchUserInput) (*models.FetchUserOutput, error) {
		userOutput, err := userHandler.FetchUser(ctx, input)
=======
>>>>>>> adc89c7 (added db layer, handler with appropriate delegation to create supabase acc, sql file, and utils file for reading sql file)
		if err != nil {
			return nil, err
		}
		return userOutput, nil
	})
}
