package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/login"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetupLoginRoutes(api huma.API, repository *data.Repository, config *config.Config) {
	loginHandler := login.NewHandler(repository.User)
	huma.Register(api, huma.Operation{
		OperationID: "login-user",
		Method:      http.MethodPost,
		Path:        "/user/login",
		Description: "user login",
		Tags:        []string{"Login"},
	}, func(ctx context.Context, input *models.LoginInput) (*models.LoginResponse, error) {
		loginOutput, err := loginHandler.Login(ctx, input, config.Supabase)
		if err != nil {
			return nil, err
		}
		return loginOutput, nil
	})
}
