package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/pfp"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

// SetUpPfpRoutes registers the logged-in user's profile picture endpoints:
// presign an upload, save it, and delete it. The user always comes from the session.
func SetUpPfpRoutes(api huma.API, repository *data.Repository, verifyRole auth.RoleVerifier) {
	pfpHandler := pfp.NewHandler(repository.User, repository.Pfp)

	huma.Register(api, huma.Operation{
		OperationID: "presign-pfp-upload",
		Middlewares: huma.Middlewares{verifyRole(api, models.StudentRole, models.CounselorRole)},
		Method:      http.MethodPut,
		Path:        "/user/pfp",
		Description: "Get a presigned URL to upload a new profile picture. PUT the image to upload_url, then send s3_key to PATCH /user/pfp to save it.",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *models.PresignPfpUploadInput) (*models.PresignPfpUploadOutput, error) {
		presigned, err := pfpHandler.PresignUpload(ctx, input.Body)
		if err != nil {
			return nil, err
		}
		return &models.PresignPfpUploadOutput{Body: *presigned}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-pfp",
		Middlewares: huma.Middlewares{verifyRole(api, models.StudentRole, models.CounselorRole)},
		Method:      http.MethodPatch,
		Path:        "/user/pfp",
		Description: "Save an uploaded image as the logged-in user's profile picture. The previous picture is deleted.",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *models.UpdatePfpInput) (*models.UpdatePfpOutput, error) {
		saved, err := pfpHandler.UpdatePfp(ctx, input.Body.S3Key)
		if err != nil {
			return nil, err
		}
		return &models.UpdatePfpOutput{Body: *saved}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-pfp",
		Middlewares: huma.Middlewares{verifyRole(api, models.StudentRole, models.CounselorRole)},
		Method:      http.MethodDelete,
		Path:        "/user/pfp",
		Description: "Remove the logged-in user's profile picture.",
		Tags:        []string{"User"},
	}, func(ctx context.Context, input *struct{}) (*struct{}, error) {
		if err := pfpHandler.DeletePfp(ctx); err != nil {
			return nil, err
		}
		return nil, nil
	})
}
