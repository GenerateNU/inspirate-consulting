package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// setupMediaTestApp creates a Fiber app in test mode with a mocked media repository.
func setupMediaTestApp(mockMedia data.MediaRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		Media: mockMedia,
	}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

// Route tests for POST /media.
func TestRoute_CreateMedia(t *testing.T) {
	t.Parallel()

	// Every field of the request body is required by Huma, so send them all.
	payload := map[string]any{
		"title":          "Common App Essay Tips",
		"description":    "How to write a strong personal statement",
		"length_in_mins": 15,
		"school_year":    nil,
		"s3_key":         "video-bucket/common-app-tips.mp4",
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// The repository should receive the decoded request body
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("CreateMedia", mock.Anything, mock.MatchedBy(func(in *models.CreateMediaRequestBody) bool {
			return in.Title == "Common App Essay Tips" && in.LengthInMins == 15 && in.S3Key == "video-bucket/common-app-tips.mp4"
		})).Return(&models.Media{
			ID:           uuid.NewString(),
			Title:        "Common App Essay Tips",
			Description:  "How to write a strong personal statement",
			LengthInMins: 15,
			S3Key:        "video-bucket/common-app-tips.mp4",
		}, nil)

		app, err := setupMediaTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodPost, "/media", payload)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var created models.Media
		require.NoError(t, json.Unmarshal(respBody, &created))
		assert.Equal(t, "Common App Essay Tips", created.Title)
		assert.Equal(t, 15, created.LengthInMins)
		assert.Nil(t, created.SchoolYear)
	})

	t.Run("validation error - missing required field", func(t *testing.T) {
		t.Parallel()

		// Mock should NOT be called when Huma validation fails
		mockRepo := mocks.NewMediaRepository(t)

		app, err := setupMediaTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/media", map[string]any{
			"title": "Common App Essay Tips",
		})

		// Huma returns 422 Unprocessable Entity for schema validation failures
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate a database failure
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("CreateMedia", mock.Anything, mock.Anything).
			Return(nil, errors.New("database error"))

		app, err := setupMediaTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/media", payload)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for GET /media/{id}.
func TestRoute_GetMedia(t *testing.T) {
	t.Parallel()

	mediaID := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// The path id should be passed through to the repository
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("GetMedia", mock.Anything, mediaID).Return(&models.Media{
			ID:           mediaID,
			Title:        "Navigating the FAFSA",
			LengthInMins: 20,
			S3Key:        "video-bucket/fafsa.mp4",
		}, nil)

		app, err := setupMediaTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodGet, "/media/"+mediaID, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var fetched models.Media
		require.NoError(t, json.Unmarshal(respBody, &fetched))
		assert.Equal(t, mediaID, fetched.ID)
		assert.Equal(t, "Navigating the FAFSA", fetched.Title)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("GetMedia", mock.Anything, mediaID).
			Return(nil, errors.New("database error"))

		app, err := setupMediaTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/media/"+mediaID, nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for GET /media.
func TestRoute_ListAllMedia(t *testing.T) {
	t.Parallel()

	t.Run("success with default pagination", func(t *testing.T) {
		t.Parallel()

		// No query params should fall back to limit=20, offset=0
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("ListAllMedia", mock.Anything, 20, 0).Return([]models.Media{
			{ID: uuid.NewString(), Title: "Common App Essay Tips", LengthInMins: 15},
			{ID: uuid.NewString(), Title: "Navigating the FAFSA", LengthInMins: 20},
		}, nil)

		app, err := setupMediaTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodGet, "/media", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var allMedia []models.Media
		require.NoError(t, json.Unmarshal(respBody, &allMedia))
		require.Len(t, allMedia, 2)
		assert.Equal(t, "Common App Essay Tips", allMedia[0].Title)
	})

	t.Run("success with limit and offset", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("ListAllMedia", mock.Anything, 1, 1).Return([]models.Media{
			{ID: uuid.NewString(), Title: "Navigating the FAFSA", LengthInMins: 20},
		}, nil)

		app, err := setupMediaTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodGet, "/media?limit=1&offset=1", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var allMedia []models.Media
		require.NoError(t, json.Unmarshal(respBody, &allMedia))
		require.Len(t, allMedia, 1)
		assert.Equal(t, "Navigating the FAFSA", allMedia[0].Title)
	})

	t.Run("validation error - invalid pagination", func(t *testing.T) {
		t.Parallel()

		for _, query := range []string{"?limit=0", "?limit=101", "?offset=-1"} {
			// Repository should not be called when pagination params are rejected
			mockRepo := mocks.NewMediaRepository(t)

			app, err := setupMediaTestApp(mockRepo)
			require.NoError(t, err)

			resp, _ := doRequest(t, app, http.MethodGet, "/media"+query, nil)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "query: %s", query)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("ListAllMedia", mock.Anything, 20, 0).
			Return(nil, errors.New("database error"))

		app, err := setupMediaTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/media", nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for DELETE /media/{id}.
func TestRoute_DeleteMedia(t *testing.T) {
	t.Parallel()

	mediaID := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("DeleteMedia", mock.Anything, mediaID).Return(nil)

		app, err := setupMediaTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/media/"+mediaID, nil)

		// Huma returns 204 No Content for an output with no body
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("DeleteMedia", mock.Anything, mediaID).
			Return(errors.New("database error"))

		app, err := setupMediaTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/media/"+mediaID, nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
