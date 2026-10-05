package routes

import (
	"inspirate-consulting/internal/auth"
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	video "inspirate-consulting/internal/handlers/video"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

// SetUpVideoRoutes registers the video endpoints:
// upload, get, and list.
func SetUpVideoRoutes(api huma.API, repository *data.Repository, verifyRole auth.RoleVerifier) {
	videoHandler := video.NewHandler(repository.Video)

	// Register POST /videos handler.
	huma.Register(api, huma.Operation{
		OperationID: "upload-video",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
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
		Middlewares: huma.Middlewares{verifyRole(api, models.StudentRole, models.CounselorRole)},
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
		Middlewares: huma.Middlewares{verifyRole(api, models.StudentRole, models.CounselorRole)},
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
