package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// setupMediaTestApp creates a Fiber app in test mode with mocked media and media access repositories.
func setupMediaTestApp(mockMedia data.MediaRepository, mockMediaAccess data.MediaAccessRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		Media:       mockMedia,
		MediaAccess: mockMediaAccess,
	}
	app, _, err := SetupApp(cfg, repo, nil, nil)
	return app, err
}

// Route tests for POST /media.
func TestRoute_CreateMedia(t *testing.T) {
	t.Parallel()

	// Send every required field plus one video source (s3_key).
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
			return in.Title == "Common App Essay Tips" && *in.LengthInMins == 15 && *in.S3Key == "video-bucket/common-app-tips.mp4"
		}), (*string)(nil)).Return(&models.Media{
			ID:           uuid.NewString(),
			Title:        "Common App Essay Tips",
			Description:  "How to write a strong personal statement",
			LengthInMins: ptr(15),
			S3Key:        ptr("video-bucket/common-app-tips.mp4"),
		}, nil)

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodPost, "/media", payload)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var created models.Media
		require.NoError(t, json.Unmarshal(respBody, &created))
		assert.Equal(t, "Common App Essay Tips", created.Title)
		assert.Equal(t, ptr(15), created.LengthInMins)
		assert.Nil(t, created.SchoolYear)
	})

	t.Run("youtube url is stored as a video id", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("CreateMedia", mock.Anything, mock.Anything, ptr("dQw4w9WgXcQ")).Return(&models.Media{
			ID:             uuid.NewString(),
			Title:          "FAFSA Walkthrough",
			YoutubeVideoID: ptr("dQw4w9WgXcQ"),
		}, nil)

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodPost, "/media", map[string]any{
			"title":       "FAFSA Walkthrough",
			"description": "Hosted on YouTube",
			"school_year": nil,
			"youtube_url": "https://youtu.be/dQw4w9WgXcQ",
		})
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var created models.Media
		require.NoError(t, json.Unmarshal(respBody, &created))
		assert.Equal(t, ptr("dQw4w9WgXcQ"), created.YoutubeVideoID)
	})

	t.Run("non-youtube url is rejected", func(t *testing.T) {
		t.Parallel()

		// Mock should NOT be called for an invalid URL
		mockRepo := mocks.NewMediaRepository(t)

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/media", map[string]any{
			"title":       "Sketchy",
			"description": "Not YouTube",
			"school_year": nil,
			"youtube_url": "https://evil.example.com/watch?v=dQw4w9WgXcQ",
		})
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("validation error - missing required field", func(t *testing.T) {
		t.Parallel()

		// Mock should NOT be called when Huma validation fails
		mockRepo := mocks.NewMediaRepository(t)

		app, err := setupMediaTestApp(mockRepo, nil)
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
		mockRepo.On("CreateMedia", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("database error"))

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/media", payload)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for GET /media/{id}.
func TestRoute_GetMedia(t *testing.T) {
	t.Parallel()

	mediaID := uuid.NewString()
	studentID := uuid.New()
	studentHeaders := map[string]string{"X-Student-ID": studentID.String()}
	counselorHeaders := map[string]string{"X-Counselor-ID": uuid.NewString()}
	fafsa := &models.Media{
		ID:           mediaID,
		Title:        "Navigating the FAFSA",
		LengthInMins: ptr(20),
		S3Key:        ptr("video-bucket/fafsa.mp4"),
	}

	t.Run("counselor can fetch any media", func(t *testing.T) {
		t.Parallel()

		// No access expectations: a counselor request must not hit the access check
		mockRepo := mocks.NewMediaRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GetMedia", mock.Anything, mediaID).Return(fafsa, nil)

		app, err := setupMediaTestApp(mockRepo, mockAccess)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodGet, "/media/"+mediaID, nil, counselorHeaders)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var fetched models.Media
		require.NoError(t, json.Unmarshal(respBody, &fetched))
		assert.Equal(t, mediaID, fetched.ID)
		assert.Equal(t, "Navigating the FAFSA", fetched.Title)
	})

	t.Run("student with access can fetch the media", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockAccess.On("HasMediaAccess", mock.Anything, studentID.String(), mediaID).Return(true, nil)
		mockRepo.On("GetMedia", mock.Anything, mediaID).Return(fafsa, nil)

		app, err := setupMediaTestApp(mockRepo, mockAccess)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/media/"+mediaID, nil, studentHeaders)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("student without access gets 404", func(t *testing.T) {
		t.Parallel()

		// No GetMedia expectation: the media must never be loaded for this student
		mockRepo := mocks.NewMediaRepository(t)
		mockAccess := mocks.NewMediaAccessRepository(t)
		mockAccess.On("HasMediaAccess", mock.Anything, studentID.String(), mediaID).Return(false, nil)

		app, err := setupMediaTestApp(mockRepo, mockAccess)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodGet, "/media/"+mediaID, nil, studentHeaders)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.NotContains(t, string(respBody), "fafsa.mp4")
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("GetMedia", mock.Anything, mediaID).
			Return(nil, errors.New("database error"))

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/media/"+mediaID, nil, counselorHeaders)
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
		mockRepo.On("ListAllMedia", mock.Anything, "", mock.Anything, 20, 0).Return([]models.Media{
			{ID: uuid.NewString(), Title: "Common App Essay Tips", LengthInMins: ptr(15)},
			{ID: uuid.NewString(), Title: "Navigating the FAFSA", LengthInMins: ptr(20)},
		}, nil)

		app, err := setupMediaTestApp(mockRepo, nil)
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
		mockRepo.On("ListAllMedia", mock.Anything, "", mock.Anything, 1, 1).Return([]models.Media{
			{ID: uuid.NewString(), Title: "Navigating the FAFSA", LengthInMins: ptr(20)},
		}, nil)

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodGet, "/media?limit=1&offset=1", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var allMedia []models.Media
		require.NoError(t, json.Unmarshal(respBody, &allMedia))
		require.Len(t, allMedia, 1)
		assert.Equal(t, "Navigating the FAFSA", allMedia[0].Title)
	})

	t.Run("search query is passed through", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("ListAllMedia", mock.Anything, "fafsa aid", mock.Anything, 20, 0).Return([]models.Media{
			{ID: uuid.NewString(), Title: "Navigating the FAFSA", LengthInMins: ptr(20)},
		}, nil)

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/media?q=fafsa%20aid", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("validation error - invalid pagination", func(t *testing.T) {
		t.Parallel()

		for _, query := range []string{"?limit=0", "?limit=101", "?offset=-1"} {
			// Repository should not be called when pagination params are rejected
			mockRepo := mocks.NewMediaRepository(t)

			app, err := setupMediaTestApp(mockRepo, nil)
			require.NoError(t, err)

			resp, _ := doRequest(t, app, http.MethodGet, "/media"+query, nil)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "query: %s", query)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("ListAllMedia", mock.Anything, "", mock.Anything, 20, 0).
			Return(nil, errors.New("database error"))

		app, err := setupMediaTestApp(mockRepo, nil)
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

		app, err := setupMediaTestApp(mockRepo, nil)
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

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/media/"+mediaID, nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for PATCH /media/{id}.
func TestRoute_UpdateMedia(t *testing.T) {
	t.Parallel()

	mediaID := uuid.NewString()
	payload := map[string]any{
		"title": "Navigating the FAFSA (2026)",
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// The path id and body title should be passed through to the repository
		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("UpdateMedia", mock.Anything, mediaID, models.UpdateMediaFields{Title: ptr("Navigating the FAFSA (2026)")}).Return(&models.Media{
			ID:           mediaID,
			Title:        "Navigating the FAFSA (2026)",
			LengthInMins: ptr(20),
			S3Key:        ptr("video-bucket/fafsa.mp4"),
		}, nil)

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodPatch, "/media/"+mediaID, payload)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updated models.Media
		require.NoError(t, json.Unmarshal(respBody, &updated))
		assert.Equal(t, mediaID, updated.ID)
		assert.Equal(t, "Navigating the FAFSA (2026)", updated.Title)
	})

	t.Run("description and youtube video are passed through", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("UpdateMedia", mock.Anything, mediaID, models.UpdateMediaFields{
			Description:    ptr("Now on YouTube"),
			YoutubeVideoID: ptr("dQw4w9WgXcQ"),
		}).Return(&models.Media{ID: mediaID, Description: "Now on YouTube", YoutubeVideoID: ptr("dQw4w9WgXcQ")}, nil)

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodPatch, "/media/"+mediaID, map[string]any{
			"description": "Now on YouTube",
			"youtube_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		})
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updated models.Media
		require.NoError(t, json.Unmarshal(respBody, &updated))
		assert.Equal(t, ptr("dQw4w9WgXcQ"), updated.YoutubeVideoID)
	})

	t.Run("empty body is rejected", func(t *testing.T) {
		t.Parallel()

		// Repository must not be called
		mockRepo := mocks.NewMediaRepository(t)

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPatch, "/media/"+mediaID, map[string]any{})
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("UpdateMedia", mock.Anything, mediaID, models.UpdateMediaFields{Title: ptr("Navigating the FAFSA (2026)")}).
			Return(nil, errs.NotFound("media", "id", mediaID))

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPatch, "/media/"+mediaID, payload)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaRepository(t)
		mockRepo.On("UpdateMedia", mock.Anything, mediaID, models.UpdateMediaFields{Title: ptr("Navigating the FAFSA (2026)")}).
			Return(nil, errors.New("database error"))

		app, err := setupMediaTestApp(mockRepo, nil)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPatch, "/media/"+mediaID, payload)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
