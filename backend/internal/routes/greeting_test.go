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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// setupTestApp creates a Fiber app in test mode with a mocked repository.
func setupTestApp(mockGreeting data.GreetingRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		Greeting: mockGreeting,
	}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

// Route tests verifying HTTP request handling, Huma validation, and responses.
func TestRoute_CreateGreeting(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// Mock the repository call expected for valid input
		mockRepo := mocks.NewGreetingRepository(t)
		mockRepo.On("CreateGreeting", mock.Anything, mock.MatchedBy(func(in models.CreateGreetingInput) bool {
			return in.Body.Name == "World"
		})).Return(&models.CreateGreetingOutput{
			Body: models.GreetingMessageBody{
				Message: "Hello, World!",
			},
		}, nil)

		app, err := setupTestApp(mockRepo)
		require.NoError(t, err)

		// Send valid JSON payload
		payload := map[string]string{
			"name": "World",
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/greeting", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify 200 OK status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify returned JSON body
		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.GreetingMessageBody
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Equal(t, "Hello, World!", output.Message)
	})

	t.Run("validation error - name exceeds maxLength", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewGreetingRepository(t)
		// Mock should NOT be called when Huma validation fails

		app, err := setupTestApp(mockRepo)
		require.NoError(t, err)

		// Send invalid payload exceeding 30 characters
		payload := map[string]string{
			"name": strings.Repeat("a", 31), // maxLength constraint is 30
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/greeting", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Huma returns 422 Unprocessable Entity for schema validation failures
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}
