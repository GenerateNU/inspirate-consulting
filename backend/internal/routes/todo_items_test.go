package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

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

// setupTodoItemTestApp creates a Fiber app in test mode with a mocked todo item repository.
func setupTodoItemTestApp(mockTodoItem data.TodoItemRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		TodoItem: mockTodoItem,
	}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

// doRequest sends a request through the app and returns the response with its body read.
func doRequest(t *testing.T, app *fiber.App, method, path string, body any) (*http.Response, []byte) {
	t.Helper()

	var reader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, path, reader)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, respBody
}

// Route tests for POST /todo-items.
func TestRoute_CreateTodoItem(t *testing.T) {
	t.Parallel()

	//studentID := uuid.NewString()
	//userID := uuid.NewString()

	// Every field of the request body is required by Huma, so send them all.
	payload := map[string]any{
		"todo_description": "Finish the Common App essay",
		"completed_at":     nil,
		"deadline":         nil,
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expectedStudentID := auth.GetStudentID(context.Background())
		expectedUserID := auth.GetUserID(context.Background())

		// The repository should receive the decoded request body
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("CreateTodoItem", mock.Anything, mock.MatchedBy(func(in *models.CreateTodoItemRequestBody) bool {
			return in.StudentID == expectedStudentID && in.TodoDescription == "Finish the Common App essay"
		})).Return(&models.TodoItem{
			ID:              uuid.NewString(),
			StudentID:       expectedStudentID,
			UserID:          expectedUserID,
			TodoDescription: "Finish the Common App essay",
		}, nil)

		app, err := setupTodoItemTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodPost, "/todo-items", payload)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var created models.TodoItem
		require.NoError(t, json.Unmarshal(respBody, &created))
		assert.Equal(t, expectedStudentID, created.StudentID)
		assert.Equal(t, "Finish the Common App essay", created.TodoDescription)
		assert.Nil(t, created.CompletedAt)
	})

	t.Run("validation error - missing required field", func(t *testing.T) {
		t.Parallel()

		// Mock should NOT be called when Huma validation fails
		mockRepo := mocks.NewTodoItemRepository(t)

		app, err := setupTodoItemTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/todo-items", map[string]any{
			"deadline" : nil,
		})

		// Huma returns 422 Unprocessable Entity for schema validation failures
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		// Simulate a database failure
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("CreateTodoItem", mock.Anything, mock.Anything).
			Return(nil, errors.New("database error"))

		app, err := setupTodoItemTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/todo-items", payload)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for GET /todo-items.
func TestRoute_GetTodoItemsByStudent(t *testing.T) {
	t.Parallel()

	// The route reads the student from the request context, not the URL
	studentID := auth.GetStudentID(context.Background())

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemsByStudent", mock.Anything, studentID).Return([]models.TodoItem{
			{ID: uuid.NewString(), StudentID: studentID, TodoDescription: "Request transcripts"},
			{ID: uuid.NewString(), StudentID: studentID, TodoDescription: "Draft personal statement"},
		}, nil)

		app, err := setupTodoItemTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodGet, "/todo-items", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var items []models.TodoItem
		require.NoError(t, json.Unmarshal(respBody, &items))
		require.Len(t, items, 2)
		assert.Equal(t, "Request transcripts", items[0].TodoDescription)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("GetTodoItemsByStudent", mock.Anything, studentID).
			Return(nil, errors.New("database error"))

		app, err := setupTodoItemTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/todo-items", nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for PATCH /todo-items/{id}.
func TestRoute_UpdateTodoItemCompletedAt(t *testing.T) {
	t.Parallel()

	itemID := uuid.NewString()
	payload := map[string]any{
		"todo_description": "Finish the Common App essay",
		"deadline":         nil,
	}

	t.Run("marking completed sets a timestamp", func(t *testing.T) {
		t.Parallel()

		completedAt := time.Now()

		// completed: true must reach the repository as a non-nil timestamp
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("UpdateTodoItemCompletedAt", mock.Anything, itemID, mock.MatchedBy(func(at *time.Time) bool {
			return at != nil
		})).Return(&models.TodoItem{
			ID:              itemID,
			TodoDescription: "Request transcripts",
			CompletedAt:     &completedAt,
		}, nil)

		app, err := setupTodoItemTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodPatch, "/todo-items/"+itemID, map[string]any{
			"completed": true,
		})
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updated models.TodoItem
		require.NoError(t, json.Unmarshal(respBody, &updated))
		assert.Equal(t, itemID, updated.ID)
		assert.NotNil(t, updated.CompletedAt)
	})

	t.Run("marking incomplete clears the timestamp", func(t *testing.T) {
		t.Parallel()

		// completed: false must reach the repository as a nil timestamp
		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("UpdateTodoItemCompletedAt", mock.Anything, itemID, (*time.Time)(nil)).
			Return(&models.TodoItem{ID: itemID, TodoDescription: "Request transcripts"}, nil)

		app, err := setupTodoItemTestApp(mockRepo)
		require.NoError(t, err)

		resp, respBody := doRequest(t, app, http.MethodPatch, "/todo-items/"+itemID, map[string]any{
			"completed": false,
		})
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updated models.TodoItem
		require.NoError(t, json.Unmarshal(respBody, &updated))
		assert.Nil(t, updated.CompletedAt)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewTodoItemRepository(t)
		mockRepo.On("CreateTodoItem", mock.Anything, mock.Anything).
			Return(nil, errors.New("database error"))

		app, err := setupTodoItemTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/todo-items", payload)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
