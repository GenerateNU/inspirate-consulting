package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// creates a Fiber app in test mode with a mocked video repository.
func setupTestAppWithVideo(mockVideo data.VideoRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		Video: mockVideo,
	}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

func TestRoute_UploadVideo(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("PresignUpload", mock.Anything, "common-app-tips.mp4").Return(&models.PresignUploadResponse{
			S3Key:     "videos/uuid-common-app-tips.mp4",
			UploadURL: "https://bucket.s3.amazonaws.com/videos/uuid-common-app-tips.mp4?X-Amz-...",
		}, nil)

		app, err := setupTestAppWithVideo(mockRepo)
		require.NoError(t, err)

		payload := map[string]string{"original_filename": "common-app-tips.mp4"}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/videos", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.PresignUploadResponse
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Equal(t, "videos/uuid-common-app-tips.mp4", output.S3Key)
	})

	t.Run("validation error : missing original_filename", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		// Mock should not be called when Huma validation fails.

		app, err := setupTestAppWithVideo(mockRepo)
		require.NoError(t, err)

		bodyBytes, err := json.Marshal(map[string]string{})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/videos", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRoute_GetVideo(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("GetVideo", mock.Anything, "videos/uuid-common-app-tips.mp4").Return(&models.Video{
			S3Key:       "videos/uuid-common-app-tips.mp4",
			DownloadURL: "https://bucket.s3.amazonaws.com/videos/uuid-common-app-tips.mp4?X-Amz-...",
		}, nil)

		app, err := setupTestAppWithVideo(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/videos/lookup?s3_key=videos/uuid-common-app-tips.mp4", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.Video
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Equal(t, "videos/uuid-common-app-tips.mp4", output.S3Key)
	})
}

func TestRoute_ListVideos(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("ListVideos", mock.Anything).Return([]models.Video{
			{S3Key: "videos/uuid-alpha.mp4", DownloadURL: "https://bucket.s3.amazonaws.com/videos/uuid-alpha.mp4?X-Amz-..."},
			{S3Key: "videos/uuid-beta.mp4", DownloadURL: "https://bucket.s3.amazonaws.com/videos/uuid-beta.mp4?X-Amz-..."},
		}, nil)

		app, err := setupTestAppWithVideo(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/videos", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.Video
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Len(t, output, 2)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewVideoRepository(t)
		mockRepo.On("ListVideos", mock.Anything).Return([]models.Video{}, nil)

		app, err := setupTestAppWithVideo(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/videos", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.Video
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Len(t, output, 0)
	})
}
