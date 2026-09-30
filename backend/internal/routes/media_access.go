package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	mediaaccess "inspirate-consulting/internal/handlers/media_access"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpMediaAccessRoutes(api huma.API, repository *data.Repository){
	mediaAccessHandler := mediaaccess.NewHandler(repository.MediaAccess)

	//Register POST /media-access handler
	huma.Register(api, huma.Operation{
		OperationID: "create-media-access",
		Method:      http.MethodPost,
		Path:        "/media-access",
		Description: "Allow a counselor to give a student access to a video",
		Tags:        []string{"Media Access"},
	}, func(ctx context.Context, input *models.GrantMediaAccessInput) (*models.GrantMediaAccessOutput, error) {
		createdAccess, err := mediaAccessHandler.GrantMediaAccess(ctx, &input.Body)
		if(err != nil){
			return nil, err
		}

		return &models.GrantMediaAccessOutput{Body: *createdAccess}, nil
	})

	//Register DELETE media-access/{id} handler
	huma.Register(api, huma.Operation{
		OperationID: "delete-media-access",
		Method:      http.MethodDelete,
		Path:        "/media-access/{id}",
		Description: "Allow a counselor to revoke a student access to a video",
		Tags:        []string{"Media Access"},
	}, func(ctx context.Context, input *models.RevokeMediaAccessInput)(*struct{}, error){
		err := mediaAccessHandler.RevokeMediaAccess(ctx, input.ID)
		if(err!=nil){
			return nil, err
		}

		return nil, nil
	})

	//Register GET /media-access handler
	huma.Register(api, huma.Operation{
		OperationID: "get-media-access",
		Method:      http.MethodGet,
		Path:        "/media-access",
		Description: "Allow students to fetch only videos they have access to",
		Tags:        []string{"Media Access"},
	}, func(ctx context.Context, input *models.ListAccessibleMediaInput)(*models.ListAccessibleMediaOutput, error){
		fetchedAccessibleMedia, err := mediaAccessHandler.ListAccessibleMedia(ctx)
		if(err!=nil){
			return nil, err
		}

		return &models.ListAccessibleMediaOutput{Body: fetchedAccessibleMedia}, nil
	})

}