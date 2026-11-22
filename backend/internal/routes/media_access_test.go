package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

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

// setupMediaAccessTestApp creates a Fiber app in test mode with a mocked media access repository.
func setupMediaAccessTestApp(mockMediaAccess data.MediaAccessRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		MediaAccess: mockMediaAccess,
	}
	app, _, err := SetupApp(cfg, repo, nil, nil)
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
	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GetMediaAccessStudentID", mock.Anything, accessID).Return(studentID, nil)
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
		mockRepo.On("GetMediaAccessStudentID", mock.Anything, accessID).Return(studentID, nil)
		mockRepo.On("RevokeMediaAccess", mock.Anything, accessID).
			Return(errors.New("database error"))

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/media-access/"+accessID, nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		// RevokeMediaAccess should NOT be called when the media access row does not exist
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("GetMediaAccessStudentID", mock.Anything, accessID).
			Return(uuid.Nil, errs.NotFound("media_access", "id", accessID))

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/media-access/"+accessID, nil)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// Route tests for GET /media-access.
func TestRoute_ListAccessibleMedia(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	idHeaders := map[string]string{
		"X-Student-ID": studentID.String(),
	}
	dueDate := time.Date(2026, 12, 1, 23, 59, 0, 0, time.UTC)

	t.Run("success with default pagination returns progress and due date", func(t *testing.T) {
		t.Parallel()

		// No query params should fall back to q="", limit=20, offset=0
		startedAt := time.Date(2026, 11, 20, 15, 0, 0, 0, time.UTC)
		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("ListAccessibleMedia", mock.Anything, studentID.String(), "", mock.Anything, 20, 0).Return([]models.AccessibleMedia{
			{
				Media:     models.Media{ID: uuid.NewString(), Title: "Common App Essay Tips", LengthInMins: ptr(15)},
				StartedAt: &startedAt,
				DueDate:   &dueDate,
			},
			{Media: models.Media{ID: uuid.NewString(), Title: "Navigating the FAFSA", LengthInMins: ptr(20)}},
		}, nil)

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodGet, "/media-access", nil, idHeaders)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var accessible []models.AccessibleMedia
		require.NoError(t, json.Unmarshal(respBody, &accessible))
		require.Len(t, accessible, 2)
		assert.Equal(t, "Common App Essay Tips", accessible[0].Title)
		require.NotNil(t, accessible[0].DueDate)
		assert.True(t, dueDate.Equal(*accessible[0].DueDate))
		assert.NotNil(t, accessible[0].StartedAt)
		assert.Nil(t, accessible[0].CompletedAt)
		assert.Nil(t, accessible[1].DueDate)

		// progress fields are flattened next to the media fields, not nested
		assert.Contains(t, string(respBody), `"due_date"`)
		assert.NotContains(t, string(respBody), `"Media"`)
	})

	t.Run("success with search, limit, and offset", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("ListAccessibleMedia", mock.Anything, studentID.String(), "fafsa", mock.Anything, 1, 1).Return([]models.AccessibleMedia{
			{Media: models.Media{ID: uuid.NewString(), Title: "Navigating the FAFSA", LengthInMins: ptr(20)}},
		}, nil)

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodGet, "/media-access?q=fafsa&limit=1&offset=1", nil, idHeaders)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var accessible []models.AccessibleMedia
		require.NoError(t, json.Unmarshal(respBody, &accessible))
		require.Len(t, accessible, 1)
		assert.Equal(t, "Navigating the FAFSA", accessible[0].Title)
	})

	t.Run("counselor is forbidden", func(t *testing.T) {
		t.Parallel()

		// A counselor has no student profile, so there is no "own" media list to return
		mockRepo := mocks.NewMediaAccessRepository(t)

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/media-access", nil, map[string]string{"X-Counselor-ID": uuid.NewString()})
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("validation error - invalid pagination", func(t *testing.T) {
		t.Parallel()

		for _, query := range []string{"?limit=0", "?limit=101", "?offset=-1"} {
			// Repository should not be called when pagination params are rejected
			mockRepo := mocks.NewMediaAccessRepository(t)

			app, err := setupMediaAccessTestApp(mockRepo)
			require.NoError(t, err)

			resp, _ := doRequest(t, app, http.MethodGet, "/media-access"+query, nil, idHeaders)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "query: %s", query)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("ListAccessibleMedia", mock.Anything, mock.Anything, "", mock.Anything, 20, 0).
			Return(nil, errors.New("database error"))

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/media-access", nil, idHeaders)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for POST /media/{id}/progress.
func TestRoute_RecordMediaProgress(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()
	mediaID := uuid.NewString()
	idHeaders := map[string]string{
		"X-Student-ID": studentID.String(),
	}
	now := time.Now().UTC()

	t.Run("started", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RecordMediaProgress", mock.Anything, studentID.String(), mediaID, false).Return(&models.MediaAccess{
			ID: uuid.NewString(), StudentID: studentID.String(), MediaID: mediaID, StartedAt: &now,
		}, nil)

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodPost, "/media/"+mediaID+"/progress", map[string]any{"status": "started"}, idHeaders)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var access models.MediaAccess
		require.NoError(t, json.Unmarshal(respBody, &access))
		assert.NotNil(t, access.StartedAt)
		assert.Nil(t, access.CompletedAt)
	})

	t.Run("completed", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RecordMediaProgress", mock.Anything, studentID.String(), mediaID, true).Return(&models.MediaAccess{
			ID: uuid.NewString(), StudentID: studentID.String(), MediaID: mediaID, StartedAt: &now, CompletedAt: &now,
		}, nil)

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/media/"+mediaID+"/progress", map[string]any{"status": "completed"}, idHeaders)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("student without access gets 404", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)
		mockRepo.On("RecordMediaProgress", mock.Anything, studentID.String(), mediaID, false).
			Return(nil, errs.NotFound("media", "id", mediaID))

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/media/"+mediaID+"/progress", map[string]any{"status": "started"}, idHeaders)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("validation error - unknown or missing status", func(t *testing.T) {
		t.Parallel()

		for _, body := range []map[string]any{{"status": "paused"}, {}} {
			// Repository should not be called when the body is rejected
			mockRepo := mocks.NewMediaAccessRepository(t)

			app, err := setupMediaAccessTestApp(mockRepo)
			require.NoError(t, err)

			resp, _ := doRequest(t, app, http.MethodPost, "/media/"+mediaID+"/progress", body, idHeaders)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, "body: %v", body)
		}
	})

	t.Run("counselor is forbidden", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewMediaAccessRepository(t)

		app, err := setupMediaAccessTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/media/"+mediaID+"/progress", map[string]any{"status": "started"},
			map[string]string{"X-Counselor-ID": uuid.NewString()})
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}
