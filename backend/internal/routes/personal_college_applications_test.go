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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func intPtr(i int) *int { return &i }

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
		GlobalCollege:              mockGlobalCollegeRepo,
	}
	app, _, err := SetupApp(cfg, repo, nil)
	return app, err
}

func TestRoute_CreatePersonalCollegeApplication(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		edDeadline := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
		college := &models.GlobalCollege{ID: 1, SchoolName: "Northeastern University", EDDeadline: &edDeadline}
		expectedApplication := &models.PersonalCollegeApplication{
			ID:              1,
			StudentID:       studentID,
			GlobalCollegeID: 1,
			ApplicationType: "ED",
			Category:        "reach",
		}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPersonalRepo.On("CreatePersonalCollegeApplication", mock.Anything, studentID, mock.MatchedBy(func(in models.CreatePersonalCollegeApplicationRequestBody) bool {
			return in.GlobalCollegeID == 1 && in.ApplicationType == "ED" && in.Category == "reach"
		})).Return(expectedApplication, nil)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPersonalRepo, mockGlobalCollegeRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"global_college_id": 1,
			"application_type":  "ED",
			"category":          "reach",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/applications", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("X-Student-ID", studentID.String())
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.PersonalCollegeApplication
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Equal(t, studentID, output.StudentID)
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
			"application_type":  "EA",
			"category":          "target",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/applications", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("X-Student-ID", studentID.String())
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

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
			"application_type":  "ED",
			"category":          "reach",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/applications", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("X-Student-ID", studentID.String())
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

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
			"application_type":  "NOT_REAL",
			"category":          "reach",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/applications", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("X-Student-ID", studentID.String())
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRoute_ListPersonalCollegeApplications(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockPersonalRepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPersonalRepo.On("ListPersonalCollegeApplicationsByStudentID", mock.Anything, studentID).
			Return([]models.PersonalCollegeApplication{
				{ID: 1, GlobalCollegeID: 1, ApplicationType: "ED", Category: "reach"},
				{ID: 2, GlobalCollegeID: 2, ApplicationType: "EA", Category: "target"},
			}, nil)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPersonalRepo, mockGlobalCollegeRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/applications", nil)
		require.NoError(t, err)
		req.Header.Set("X-Student-ID", studentID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

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
		mockPersonalRepo.On("ListPersonalCollegeApplicationsByStudentID", mock.Anything, studentID).
			Return([]models.PersonalCollegeApplication{}, nil)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPersonalRepo, mockGlobalCollegeRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/applications", nil)
		require.NoError(t, err)
		req.Header.Set("X-Student-ID", studentID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.PersonalCollegeApplication
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Len(t, output, 0)
	})
}

func TestRoute_UpdatePersonalCollegeApplication(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		deadline := time.Now().Add(24 * time.Hour)
		college := &models.GlobalCollege{ID: 1, SchoolName: "Test U", RDDeadline: &deadline}

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockGlobalCollegeRepo.On("GetGlobalCollege", mock.Anything, int64(1)).Return(college, nil)

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("UpdatePersonalCollegeApplication", mock.Anything, studentID, int64(42), mock.Anything).
			Return(&models.PersonalCollegeApplication{ID: 42, GlobalCollegeID: 1, ApplicationType: "RD", Category: "safety"}, nil)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPCARepo, mockGlobalCollegeRepo)
		require.NoError(t, err)

		payload := map[string]any{"global_college_id": 1, "application_type": "RD", "category": "safety"}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "/applications/42", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("X-Student-ID", studentID.String())
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("validation error : invalid application_type", func(t *testing.T) {
		t.Parallel()

		mockGlobalCollegeRepo := mocks.NewGlobalCollegeRepository(t)
		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPCARepo, mockGlobalCollegeRepo)
		require.NoError(t, err)

		payload := map[string]any{"global_college_id": 1, "application_type": "BOGUS", "category": "safety"}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "/applications/42", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("X-Student-ID", studentID.String())
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRoute_UpdatePersonalCollegeApplicationRank(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("UpdateApplicationRank", mock.Anything, studentID, int64(42), intPtr(1)).
			Return([]models.PersonalCollegeApplication{{ID: 42, Rank: intPtr(1)}}, nil)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPCARepo, nil)
		require.NoError(t, err)

		bodyBytes, err := json.Marshal(map[string]any{"rank": 1})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "/applications/42/rank", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("X-Student-ID", studentID.String())
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		var output []models.PersonalCollegeApplication
		require.NoError(t, json.Unmarshal(respBody, &output))
		assert.Len(t, output, 1)
	})

	t.Run("success : unrank with null", func(t *testing.T) {
		t.Parallel()

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("UpdateApplicationRank", mock.Anything, studentID, int64(42), (*int)(nil)).
			Return([]models.PersonalCollegeApplication{{ID: 42, Rank: nil}}, nil)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPCARepo, nil)
		require.NoError(t, err)

		bodyBytes, err := json.Marshal(map[string]any{"rank": nil})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, "/applications/42/rank", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("X-Student-ID", studentID.String())
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestRoute_DeletePersonalCollegeApplication(t *testing.T) {
	t.Parallel()

	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockPCARepo := mocks.NewPersonalCollegeApplicationRepository(t)
		mockPCARepo.On("DeletePersonalCollegeApplication", mock.Anything, studentID, int64(42)).Return(nil)

		app, err := setupTestAppWithPersonalCollegeApplication(mockPCARepo, nil)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodDelete, "/applications/42", nil)
		require.NoError(t, err)
		req.Header.Set("X-Student-ID", studentID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
