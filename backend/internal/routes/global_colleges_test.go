package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// setupTestAppWithGlobalCollege creates a Fiber app in test mode with a
// mocked global college repository.
func setupTestAppWithGlobalCollege(mockGlobalCollege data.GlobalCollegeRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		GlobalCollege: mockGlobalCollege,
	}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

func TestRoute_CreateGlobalCollege(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		edDeadline := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("CreateGlobalCollege", mock.Anything, mock.MatchedBy(func(in models.CreateGlobalCollegeRequestBody) bool {
			return in.SchoolName == "Northeastern University" && in.SchoolLocation == "Boston, MA"
		})).Return(&models.GlobalCollege{
			ID: 1,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
			SchoolName: "Northeastern University",
			SchoolLocation: "Boston, MA",
			EDDeadline: &edDeadline,
		}, nil)

		app, err := setupTestAppWithGlobalCollege(mockRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"school_name": "Northeastern University",
			"school_location": "Boston, MA",
			"ed_deadline": edDeadline,
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/colleges", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.GlobalCollege
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Equal(t, "Northeastern University", output.SchoolName)
		assert.Equal(t, "Boston, MA", output.SchoolLocation)
	})

	t.Run("duplicate returns conflict", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("CreateGlobalCollege", mock.Anything, mock.MatchedBy(func(in models.CreateGlobalCollegeRequestBody) bool {
			return in.SchoolName == "Duplicate University"
		})).Return(nil, errs.Conflict("global college", "school_name/school_location", "Duplicate University / Dupe City, DC"))

		app, err := setupTestAppWithGlobalCollege(mockRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"school_name": "Duplicate University",
			"school_location": "Dupe City, DC",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/colleges", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		// The repository's unique-index violation is translated to errs.Conflict,
		// which the central ErrorHandler surfaces as a real 409.
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("validation error : missing required school_name", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		// Mock should NOT be called when Huma validation fails.

		app, err := setupTestAppWithGlobalCollege(mockRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"school_location": "Boston, MA",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/colleges", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRoute_ListGlobalColleges(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("ListGlobalColleges", mock.Anything).Return([]models.GlobalCollege{
			{ID: 1, CreatedAt: createdAt, UpdatedAt: createdAt, SchoolName: "Alpha University", SchoolLocation: "Alpha City, AA"},
			{ID: 2, CreatedAt: createdAt, UpdatedAt: createdAt, SchoolName: "Beta College", SchoolLocation: "Beta Town, BB"},
		}, nil)

		app, err := setupTestAppWithGlobalCollege(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/colleges", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.GlobalCollege
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Len(t, output, 2)
		assert.Equal(t, "Alpha University", output[0].SchoolName)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("ListGlobalColleges", mock.Anything).Return([]models.GlobalCollege{}, nil)

		app, err := setupTestAppWithGlobalCollege(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/colleges", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.GlobalCollege
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Len(t, output, 0)
	})
}

func TestRoute_GetGlobalCollege(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(&models.GlobalCollege{
			ID: 1,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
			SchoolName: "Northeastern University",
			SchoolLocation: "Boston, MA",
		}, nil)

		app, err := setupTestAppWithGlobalCollege(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/colleges/1", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.GlobalCollege
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Equal(t, "Northeastern University", output.SchoolName)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGlobalCollegeRepository(t)
		mockRepo.On("GetGlobalCollege", mock.Anything, int64(999)).
			Return(nil, errs.NotFound("global college", "id", "999"))

		app, err := setupTestAppWithGlobalCollege(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/colleges/999", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		// GetGlobalCollege's repository returns an errs.NotFound (HTTPError).
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}