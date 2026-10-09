package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
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

// setupEssayTestApp creates a Fiber app in test mode with a mocked essay repository.
func setupEssayTestApp(mockEssay data.EssayRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		Essay: mockEssay,
	}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

// Route tests verifying HTTP request handling, Huma validation, and responses.
func TestRoute_GetEssaysFromStudent(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()
		collegeID := int64(42)

		// Mock the repository call expected for valid input
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssaysFromStudent", mock.Anything, studentID).Return([]models.Essays{
			{
				ID:            uuid.New(),
				StudentID:     studentID,
				Type:          "personal-statement",
				CollegeID:     &collegeID,
				LinkToContent: "https://docs.google.com/document/d/abc123",
				Status:        models.Draft,
			},
		}, nil)

		app, err := setupEssayTestApp(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/students/"+studentID.String()+"/essays", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		// Verify 200 OK status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify returned JSON body
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.EssayListBody
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.Len(t, output.Essays, 1)
		assert.Equal(t, studentID, output.Essays[0].StudentID)
	})

	t.Run("validation error - student id is not a uuid", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		// Mock should NOT be called when Huma validation fails

		app, err := setupEssayTestApp(mockRepo)
		require.NoError(t, err)

		// Send a path value that cannot be parsed as a uuid
		req, err := http.NewRequest(http.MethodGet, "/students/not-a-uuid/essays", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		// Huma returns 422 Unprocessable Entity for schema validation failures
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRoute_CreateEssay(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()

		// Mock the repository call expected for valid input. ID and Status are
		// left unset so the database applies its own defaults.
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("CreateEssay", mock.Anything, mock.MatchedBy(func(essay models.Essays) bool {
			return essay.StudentID == studentID &&
				essay.Type == "personal-statement" &&
				essay.ID == uuid.Nil &&
				essay.Status == ""
		})).Return(nil)

		app, err := setupEssayTestApp(mockRepo)
		require.NoError(t, err)

		// Send valid JSON payload
		payload := map[string]any{
			"student_id":      studentID.String(),
			"type":            "personal-statement",
			"link_to_content": "https://docs.google.com/document/d/abc123",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/essays", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		// The output envelope has no Body, so Huma responds 204 No Content
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("validation error - type exceeds maxLength", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		// Mock should NOT be called when Huma validation fails

		app, err := setupEssayTestApp(mockRepo)
		require.NoError(t, err)

		// Send invalid payload exceeding 100 characters
		payload := map[string]any{
			"student_id":      uuid.New().String(),
			"type":            strings.Repeat("a", 101),
			"link_to_content": "https://docs.google.com/document/d/abc123",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/essays", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		// Huma returns 422 Unprocessable Entity for schema validation failures
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRoute_UpdateEssayStatus(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		essayID := uuid.New()

		// Mock the repository call expected for valid input
		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On(
			"UpdateStatus",
			mock.Anything,
			essayID,
			models.Submitted,
		).Return(&models.Essays{
			ID:     essayID,
			Status: models.Submitted,
		}, nil)

		app, err := setupEssayTestApp(mockRepo)
		require.NoError(t, err)

		// Send valid JSON payload
		payload := map[string]string{
			"status": "Submitted",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPatch,
			"/essays/"+essayID.String()+"/status",
			bytes.NewReader(bodyBytes),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify the updated essay is returned in the body
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.EssayBody
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.NotNil(t, output.Essay)
		assert.Equal(t, models.Submitted, output.Essay.Status)
	})

	t.Run("validation error - status is not one of the enum values", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		// Mock should NOT be called when Huma validation fails

		app, err := setupEssayTestApp(mockRepo)
		require.NoError(t, err)

		// Send a status outside the four declared enum values
		payload := map[string]string{
			"status": "Rejected",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPatch,
			"/essays/"+uuid.New().String()+"/status",
			bytes.NewReader(bodyBytes),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		// Huma returns 422 Unprocessable Entity for schema validation failures
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRoute_GetEssaysByGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		groupID := uuid.New()

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("GetEssaysByGroup", mock.Anything, groupID).Return([]models.Essays{
			{
				ID:            uuid.New(),
				StudentID:     uuid.New(),
				Type:          "personal-statement",
				LinkToContent: "https://docs.google.com/document/d/abc123",
				Status:        models.Draft,
				EssayGroupID:  &groupID,
			},
		}, nil)

		app, err := setupEssayTestApp(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/essay-groups/"+groupID.String()+"/essays", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.EssayListBody
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.Len(t, output.Essays, 1)
		require.NotNil(t, output.Essays[0].EssayGroupID)
		assert.Equal(t, groupID, *output.Essays[0].EssayGroupID)
	})

	t.Run("validation error - group id is not a uuid", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		// Mock should NOT be called when Huma validation fails

		app, err := setupEssayTestApp(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/essay-groups/not-a-uuid/essays", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRoute_UpdateEssayGroup(t *testing.T) {
	t.Parallel()

	t.Run("moves the essay into a group", func(t *testing.T) {
		t.Parallel()

		essayID := uuid.New()
		groupID := uuid.New()

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("UpdateEssayGroup", mock.Anything, essayID, mock.MatchedBy(func(id *uuid.UUID) bool {
			return id != nil && *id == groupID
		})).Return(&models.Essays{
			ID:           essayID,
			EssayGroupID: &groupID,
		}, nil)

		app, err := setupEssayTestApp(mockRepo)
		require.NoError(t, err)

		payload := map[string]any{"essay_group_id": groupID.String()}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPatch,
			"/essays/"+essayID.String()+"/group",
			bytes.NewReader(bodyBytes),
		)
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

		var output models.EssayBody
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.NotNil(t, output.Essay)
		require.NotNil(t, output.Essay.EssayGroupID)
		assert.Equal(t, groupID, *output.Essay.EssayGroupID)
	})

	// Omitting essay_group_id entirely is how a client removes an essay from its group.
	t.Run("removes the essay from its group", func(t *testing.T) {
		t.Parallel()

		essayID := uuid.New()

		mockRepo := mocks.NewEssayRepository(t)
		mockRepo.On("UpdateEssayGroup", mock.Anything, essayID, (*uuid.UUID)(nil)).Return(&models.Essays{
			ID:           essayID,
			EssayGroupID: nil,
		}, nil)

		app, err := setupEssayTestApp(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPatch,
			"/essays/"+essayID.String()+"/group",
			strings.NewReader(`{}`),
		)
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

		var output models.EssayBody
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.NotNil(t, output.Essay)
		assert.Nil(t, output.Essay.EssayGroupID)
	})

	t.Run("validation error - essay group id is not a uuid", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayRepository(t)
		// Mock should NOT be called when Huma validation fails

		app, err := setupEssayTestApp(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPatch,
			"/essays/"+uuid.New().String()+"/group",
			strings.NewReader(`{"essay_group_id": "not-a-uuid"}`),
		)
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
