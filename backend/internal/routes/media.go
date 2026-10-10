package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	media "inspirate-consulting/internal/handlers/media"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpMediaRoutes(api huma.API, repository *data.Repository) {
	mediaHandler := media.NewHandler(repository.Media)

	//Register POST /media handler
	huma.Register(api, huma.Operation{
		OperationID: "create-media",
		Method:      http.MethodPost,
		Path:        "/media",
		Description: "Create a media/video",
		Tags:        []string{"Media"},
	}, func(ctx context.Context, input *models.CreateMediaInput) (*models.CreateMediaOutput, error) {
		created, err := mediaHandler.CreateMedia(ctx, &input.Body)
		if err != nil {
			return nil, err
		}

		return &models.CreateMediaOutput{Body: *created}, nil
	})

	//Register GET /media/{id} handler
	huma.Register(api, huma.Operation{
		OperationID: "get-media",
		Method:      http.MethodGet,
		Path:        "/media/{id}",
		Description: "Get a media/video",
		Tags:        []string{"Media"},
	}, func(ctx context.Context, input *models.GetMediaInput) (*models.GetMediaOutput, error) {
		fetchedMedia, err := mediaHandler.GetMedia(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return &models.GetMediaOutput{Body: *fetchedMedia}, nil
	})

	//Register GET /media handler
	huma.Register(api, huma.Operation{
		OperationID: "list-media",
		Method:      http.MethodGet,
		Path:        "/media",
		Description: "Fetch all videos, through counselor view",
		Tags:        []string{"Media"},
	}, func(ctx context.Context, input *models.ListAllMediaInput) (*models.ListAllMediaOutput, error) {
		allMedia, err := mediaHandler.ListAllMedia(ctx, input.Limit, input.Offset)
		if err != nil {
			return nil, err
		}
		return &models.ListAllMediaOutput{Body: allMedia}, nil
	})
	//Register DELETE media/{id} handler
	huma.Register(api, huma.Operation{
		OperationID: "delete-media",
		Method:      http.MethodDelete,
		Path:        "/media/{id}",
		Description: "Delete a video, through counselor view",
		Tags:        []string{"Media"},
	}, func(ctx context.Context, input *models.DeleteMediaInput) (*struct{}, error) {
		err := mediaHandler.DeleteMedia(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

}
