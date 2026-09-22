package routes

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func setupTestAppWithStudent(mockRepo data.EssayReviewRepository) (*fiber.App, error) {
	cfg := config.Config{TestMode: true}
	repo := &data.Repository{EssayReview: mockRepo}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

func TestRoute_SetStudentReviewBalance(t *testing.T) {
	t.Parallel()

	t.Run("success records the delta and returns the student", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(&models.Student{ID: studentID, Year: models.Senior, GPA: 4, ReviewBalance: 2}, nil)
		mockRepo.On("SetStudentBalance", mock.Anything, mock.Anything, studentID, 7).Return(nil)
		mockRepo.On("InsertAdjustment", mock.Anything, mock.Anything, studentID, 5).Return(nil)

		app, err := setupTestAppWithStudent(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodPatch, "/students/"+studentID.String()+"/review-balance",
			map[string]any{"review_balance": 7})

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var out models.Student
		require.NoError(t, json.Unmarshal(body, &out))
		assert.Equal(t, 7, out.ReviewBalance)
		assert.Equal(t, studentID, out.ID)
	})

	t.Run("setting the same value writes nothing", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(&models.Student{ID: studentID, ReviewBalance: 4}, nil)

		app, err := setupTestAppWithStudent(mockRepo)
		require.NoError(t, err)

		resp, _ := doJSON(t, app, http.MethodPatch, "/students/"+studentID.String()+"/review-balance",
			map[string]any{"review_balance": 4})

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockRepo.AssertNotCalled(t, "InsertAdjustment", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("negative balance is rejected by validation", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayReviewRepository(t)
		// The handler is never reached.

		app, err := setupTestAppWithStudent(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodPatch, "/students/"+uuid.New().String()+"/review-balance",
			map[string]any{"review_balance": -1})

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Contains(t, string(body), "expected number >= 0")
	})

	t.Run("unknown student returns not found", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(nil, errs.NotFound("student", "id", studentID.String()))

		app, err := setupTestAppWithStudent(mockRepo)
		require.NoError(t, err)

		resp, _ := doJSON(t, app, http.MethodPatch, "/students/"+studentID.String()+"/review-balance",
			map[string]any{"review_balance": 3})

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
