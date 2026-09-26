package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	video "inspirate-consulting/internal/handlers/video"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

// SetUpVideoRoutes registers the video endpoints:
// upload, get, and list.
func SetUpVideoRoutes(api huma.API, repository *data.Repository) {
	videoHandler := video.NewHandler(repository.Video)

	// Register POST /videos handler.
	huma.Register(api, huma.Operation{
		OperationID: "upload-video",
		Method:      http.MethodPost,
		Path:        "/videos",
		Description: "Upload a new video to the S3 bucket.",
		Tags:        []string{"Videos"},
	}, func(ctx context.Context, input *models.PresignUploadInput) (*models.PresignUploadOutput, error) {
		created, err := videoHandler.UploadVideo(ctx, input.Body.OriginalFilename)
		if err != nil {
			return nil, err
		}
		return &models.PresignUploadOutput{Body: *created}, nil
	})

	// Register GET /videos/lookup handler for retrieving a video by its S3 key.
	huma.Register(api, huma.Operation{
		OperationID: "get-video",
		Method:      http.MethodGet,
		Path:        "/videos/lookup",
		Description: "Get a video by its S3 key.",
		Tags:        []string{"Videos"},
	}, func(ctx context.Context, input *models.GetVideoInput) (*models.GetVideoOutput, error) {
		video, err := videoHandler.GetVideo(ctx, input.S3Key)
		if err != nil {
			return nil, err
		}
		return &models.GetVideoOutput{Body: *video}, nil
	})

	// Register GET /videos handler for listing all videos in the S3 bucket.
	huma.Register(api, huma.Operation{
		OperationID: "list-videos",
		Method:      http.MethodGet,
		Path:        "/videos",
		Description: "List all videos in the S3 bucket.",
		Tags:        []string{"Videos"},
	}, func(ctx context.Context, input *models.ListVideosInput) (*models.ListVideosOutput, error) {
		videos, err := videoHandler.ListVideos(ctx)
		if err != nil {
			return nil, err
		}
		return &models.ListVideosOutput{Body: videos}, nil
	})
}
