package routes

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/data/aws/sqsclient"
	resetpassword "inspirate-consulting/internal/handlers/reset_password"
	"inspirate-consulting/internal/models"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func SetupResetPasswordRoutes(api huma.API, repository *data.Repository, config *config.Config, verifyRole auth.RoleVerifier, accountSQSClient *sqsclient.SQSClient) {
	resetpasswordHandler := resetpassword.NewHandler(repository.User, accountSQSClient)
	huma.Register(api, huma.Operation{
		OperationID: "reset-user-password",
		Method:      http.MethodPatch,
		Path:        "/user/reset-password",
		Description: "reset user password",
		Tags:        []string{"ResetPass"},
	}, func(ctx context.Context, input *models.ResetPasswordInput) (*models.ResetPasswordResponse, error) {
		resetpasswordOutput, err := resetpasswordHandler.ResetPassword(ctx, input, config.Supabase)
		if err != nil {
			return nil, err
		}
		return resetpasswordOutput, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "force-password-reset",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
		Method:      http.MethodPost,
		Path:        "/user/{supabase_id}/force-password-reset",
		Description: "Force a user to reset their password on next login",
		Tags:        []string{"ResetPass"},
	}, func(ctx context.Context, input *models.ForcePasswordResetInput) (*models.ForcePasswordResetOutput, error) {
		return resetpasswordHandler.ForcePasswordReset(ctx, input)
	})

	huma.Register(api, huma.Operation{
		OperationID: "send-temp-password",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
		Method:      http.MethodPost,
		Path:        "/user/{id}/send-temp-password",
		Description: "Generate a new temporary password for a user and email it to them",
		Tags:        []string{"ResetPass"},
	}, func(ctx context.Context, input *models.SendTempPasswordInput) (*models.SendTempPasswordOutput, error) {
		return resetpasswordHandler.SendTempPassword(ctx, input, config.Supabase)
	})
}
