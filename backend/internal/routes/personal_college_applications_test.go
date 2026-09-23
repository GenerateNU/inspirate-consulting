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

// setupTestAppWithPersonalCollegeApplication creates a Fiber app in test
// mode with mocked personal college application and global college
// repositories (the handler depends on both).
func setupTestAppWithPersonalCollegeApplication(
	mockPersonalRepo data.PersonalCollegeApplicationRepository,
	mockGlobalCollegeRepo data.GlobalCollegeRepository,
) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		PersonalCollegeApplication: mockPersonalRepo,
		GlobalCollege:      mockGlobalCollegeRepo,
	}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

func TestRoute_CreatePersonalCollegeApplication(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		edDeadline := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
		college := &models.GlobalCollege{ID: 1, SchoolName: "Northeastern University", EDDeadline: &edDeadline}
		expectedApplication := &models.PersonalCollegeApplication{
			ID: 1,
			GlobalCollegeID: 1,
			ApplicationType: "ED",
			Category: "reach",
		}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPersonalRepo.On("CreatePersonalCollegeApplication", mock.Anything, mock.Anything, mock.MatchedBy(func(in models.CreatePersonalCollegeApplicationRequestBody) bool {
			return in.GlobalCollegeID == 1 && in.ApplicationType == "ED" && in.Category == "reach"
		})).Return(expectedApplication, nil)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPersonalRepo, mockGlobalCollegeRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"global_college_id": 1,
			"application_type": "ED",
			"category": "reach",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/applications", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.PersonalCollegeApplication
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Equal(t, "ED", output.ApplicationType)
		assert.Equal(t, "reach", output.Category)
	})

	t.Run("deadline type not offered returns bad request", func(t *testing.T) {
		t.Parallel()

		// College only offers ED; EADeadline is nil.
		edDeadline := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
		college := &models.GlobalCollege{ID: 1, SchoolName: "Northeastern University", EDDeadline: &edDeadline}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		// CreatePersonalCollegeApplication will not be called.

		app, err := setupTestAppWithPersonalCollegeApplication(mockPersonalRepo, mockGlobalCollegeRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"global_college_id": 1,
			"application_type": "EA",
			"category": "target",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/applications", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("nonexistent global college returns not found", func(t *testing.T) {
		t.Parallel()

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(999)).
			Return(nil, errs.NotFound("global college", "id", "999"))

		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPersonalRepo, mockGlobalCollegeRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"global_college_id": 999,
			"application_type": "ED",
			"category": "reach",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/applications", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("validation error : invalid application_type", func(t *testing.T) {
		t.Parallel()

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		// Huma will reject the request before reaching the handler

		app, err := setupTestAppWithPersonalCollegeApplication(mockPersonalRepo, mockGlobalCollegeRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"global_college_id": 1,
			"application_type": "NOT_REAL",
			"category": "reach",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/applications", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRoute_ListPersonalCollegeApplications(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPersonalRepo.On("ListPersonalCollegeApplicationsByStudentID", mock.Anything, mock.Anything).
			Return([]models.PersonalCollegeApplication{
				{ID: 1, GlobalCollegeID: 1, ApplicationType: "ED", Category: "reach"},
				{ID: 2, GlobalCollegeID: 2, ApplicationType: "EA", Category: "target"},
			}, nil)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPersonalRepo, mockGlobalCollegeRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/applications", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.PersonalCollegeApplication
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Len(t, output, 2)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPersonalRepo.On("ListPersonalCollegeApplicationsByStudentID", mock.Anything, mock.Anything).
			Return([]models.PersonalCollegeApplication{}, nil)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPersonalRepo, mockGlobalCollegeRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/applications", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.PersonalCollegeApplication
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Len(t, output, 0)
	})
}