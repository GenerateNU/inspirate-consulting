package routes

import (
	"context"
	"fmt"
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
		fmt.Println("[route] CreateUser called")
		userOutput, err := userHandler.CreateUser(ctx, input, config)
		if err != nil {
			fmt.Println("[route] CreateUser error:", err)
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
		if err != nil {
			return nil, err
		}
		return userOutput, nil
	})
}
