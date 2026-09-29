package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"inspirate-consulting/internal/auth"
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

// setupMediaAccessTestApp creates a Fiber app in test mode with a mocked media access repository.
func setupMediaAccessTestApp(mockMediaAccess data.MediaAccessRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		MediaAccess: mockMediaAccess,
	}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

// Route tests for POST /media-access.
func TestRoute_GrantMediaAccess(t *testing.T) {
	t.Parallel()

	studentID := uuid.NewString()
	mediaID := uuid.NewString()

	// Every field of the request body is required by Huma, so send them all.
	payload := map[string]any{
		"student_id": studentID,
		"media_id":   mediaID,
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// The repository should receive the decoded request body
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GrantMediaAccess", mock.Anything, mock.MatchedBy(func(in *models.GrantMediaAccessRequestBody) bool {
			return in.StudentID == studentID && in.MediaID == mediaID
		})).Return(&models.MediaAccess{
			ID:        uuid.NewString(),
			StudentID: studentID,
			MediaID:   mediaID,
		}, nil)

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodPost, "/media-access", payload)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var created models.MediaAccess
		require.NoError(t, json.Unmarshal(respBody, &created))
		assert.Equal(t, studentID, created.StudentID)
		assert.Equal(t, mediaID, created.MediaID)
	})

	t.Run("validation error - missing required field", func(t *testing.T) {
		t.Parallel()

		// Mock should NOT be called when Huma validation fails
		mockRepo := mocks.NewMediaAccessRepository(t)

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/media-access", map[string]any{
			"student_id": studentID,
		})

		// Huma returns 422 Unprocessable Entity for schema validation failures
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate a database failure
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GrantMediaAccess", mock.Anything, mock.Anything).
			Return(nil, errors.New("database error"))

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/media-access", payload)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for DELETE /media-access/{id}.
func TestRoute_RevokeMediaAccess(t *testing.T) {
	t.Parallel()

	accessID := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RevokeMediaAccess", mock.Anything, accessID).Return(nil)

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/media-access/"+accessID, nil)

		// Huma returns 204 No Content for an output with no body
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RevokeMediaAccess", mock.Anything, accessID).
			Return(errors.New("database error"))

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/media-access/"+accessID, nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for GET /media-access.
func TestRoute_ListAccessibleMedia(t *testing.T) {
	t.Parallel()

	// The route reads the student from the request context, not the URL
	studentID := auth.GetStudentID(context.Background())

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("ListAccessibleMedia", mock.Anything, studentID).Return([]models.Media{
			{ID: uuid.NewString(), Title: "Common App Essay Tips", LengthInMins: 15},
			{ID: uuid.NewString(), Title: "Navigating the FAFSA", LengthInMins: 20},
		}, nil)

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodGet, "/media-access", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var accessible []models.Media
		require.NoError(t, json.Unmarshal(respBody, &accessible))
		require.Len(t, accessible, 2)
		assert.Equal(t, "Common App Essay Tips", accessible[0].Title)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("ListAccessibleMedia", mock.Anything, studentID).
			Return(nil, errors.New("database error"))

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/media-access", nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
