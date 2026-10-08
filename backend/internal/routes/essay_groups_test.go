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

// setupEssayGroupTestApp creates a Fiber app in test mode with a mocked essay group repository.
func setupEssayGroupTestApp(mockEssayGroup data.EssayGroupRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		EssayGroup: mockEssayGroup,
	}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

func TestRoute_CreateEssayGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()
		groupID := uuid.New()

		mockRepo := mocks.NewEssayGroupRepository(t)
		mockRepo.On("CreateEssayGroup", mock.Anything, mock.MatchedBy(func(body models.CreateEssayGroupBody) bool {
			return body.Name == "Common App" && body.StudentID == studentID
		})).Return(&models.EssayGroups{
			ID:        groupID,
			Name:      "Common App",
			StudentID: studentID,
		}, nil)

		app, err := setupEssayGroupTestApp(mockRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"name":       "Common App",
			"student_id": studentID.String(),
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/essay-groups", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify the created group is returned so the client can use its id
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.EssayGroups
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Equal(t, groupID, output.ID)
		assert.Equal(t, "Common App", output.Name)
	})

	t.Run("validation error - name exceeds maxLength", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayGroupRepository(t)
		// Mock should NOT be called when Huma validation fails

		app, err := setupEssayGroupTestApp(mockRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"name":       strings.Repeat("a", 101),
			"student_id": uuid.New().String(),
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/essay-groups", bytes.NewReader(bodyBytes))
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

	t.Run("validation error - name is empty", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayGroupRepository(t)
		// Mock should NOT be called when Huma validation fails

		app, err := setupEssayGroupTestApp(mockRepo)
		require.NoError(t, err)

		payload := map[string]any{
			"name":       "",
			"student_id": uuid.New().String(),
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/essay-groups", bytes.NewReader(bodyBytes))
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

func TestRoute_GetEssayGroupsFromStudent(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()
		description := "Essays for the Common App"

		mockRepo := mocks.NewEssayGroupRepository(t)
		mockRepo.On("ListEssayGroups", mock.Anything, studentID).Return([]models.EssayGroups{
			{
				ID:          uuid.New(),
				Name:        "Common App",
				Description: &description,
				StudentID:   studentID,
			},
		}, nil)

		app, err := setupEssayGroupTestApp(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/students/"+studentID.String()+"/essay-groups", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.EssayGroupListBody
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.Len(t, output.EssayGroups, 1)
		assert.Equal(t, studentID, output.EssayGroups[0].StudentID)
	})

	t.Run("validation error - student id is not a uuid", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayGroupRepository(t)
		// Mock should NOT be called when Huma validation fails

		app, err := setupEssayGroupTestApp(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/students/not-a-uuid/essay-groups", nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}
