package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/config"
	claimpassword "inspirate-consulting/internal/handlers/claim_password"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetupClaimPasswordRoutes(api huma.API, config *config.Config) {
	claimPasswordHandler := claimpassword.NewHandler(config.ClaimToken.Secret)

	huma.Register(api, huma.Operation{
		OperationID: "claim-password",
		Method:      http.MethodGet,
		Path:        "/claim-password",
		Description: "Exchange a claim-password token (from an account-created or password-reset email) for the temporary password it carries",
		Tags:        []string{"ClaimPassword"},
	}, func(ctx context.Context, input *models.ClaimPasswordInput) (*models.ClaimPasswordOutput, error) {
		return claimPasswordHandler.ClaimPassword(ctx, input)
	})
}
