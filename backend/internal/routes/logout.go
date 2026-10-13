package routes

import (
	"context"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/logout"
	"inspirate-consulting/internal/models"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func SetupLogoutRoutes(api huma.API, repository *data.Repository, config *config.Config) {
	logoutHandler := logout.NewHandler(repository.User)
	huma.Register(api, huma.Operation{
		OperationID: "logout-user",
		Method:      http.MethodPost,
		Path:        "/user/logout",
		Description: "user logout",
		Tags:        []string{"Logout"},
	}, func(ctx context.Context, input *models.LogoutInput) (*models.LogoutResponse, error) {
		logoutOutput, err := logoutHandler.Logout(ctx, input, config.Supabase)
		if err != nil {
			return nil, err
		}
		return logoutOutput, nil
	})
}
