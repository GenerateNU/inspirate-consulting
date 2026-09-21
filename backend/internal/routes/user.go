package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/user"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetupUserRoutes(api huma.API, repository *data.Repository) {
	userHandler := user.NewHandler(repository.User)
	huma.Register(api, huma.Operation{
		OperationID: "create-user",
		Method:      http.MethodPost,
		Path:        "/user",
		Description: "Create a user (student/counselor)",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *models.CreateUserInput) (*models.CreateUserOutput, error) {
		userOutput, err := userHandler.CreateUser(ctx, input)
		if err != nil {
			return nil, err
		}
		return userOutput, nil
	})
}
