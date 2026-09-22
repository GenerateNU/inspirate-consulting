package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	dbinterface "inspirate-consulting/internal/data/db-interface"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// setupTestAppWithEssayReview creates a Fiber app in test mode with a mocked
// essay review repository.
func setupTestAppWithEssayReview(mockRepo data.EssayReviewRepository) (*fiber.App, error) {
	cfg := config.Config{TestMode: true}
	repo := &data.Repository{EssayReview: mockRepo}
	app, _, err := SetupApp(cfg, repo)
	return app, err
}

// expectTx makes the mock run the callback the handler passes to WithTx.
func expectTx(mockRepo *mocks.EssayReviewRepository) {
	mockRepo.On("WithTx", mock.Anything, mock.Anything).
		Return(func(ctx context.Context, fn func(db dbinterface.QueryInterface) error) error {
			return fn(nil)
		})
}

func doJSON(t *testing.T, app *fiber.App, method, path string, payload any) (*http.Response, []byte) {
	t.Helper()

	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		require.NoError(t, err)
		body = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, path, body)
	require.NoError(t, err)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp, respBody
}

func TestRoute_RequestEssayReview(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		studentID, essayID := uuid.New(), uuid.New()
		status := models.ReviewStatusOpen
		created := &models.EssayReviewTransaction{
			ID: uuid.New(), StudentID: studentID, EssayID: &essayID,
			Subtotal: -1, EntryType: models.ReviewEntrySpend, Status: &status,
		}

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(&models.Student{ID: studentID, ReviewBalance: 5}, nil)
		mockRepo.On("FindOpenReviewForEssay", mock.Anything, mock.Anything, essayID).Return(nil, nil)
		mockRepo.On("AdjustStudentBalance", mock.Anything, mock.Anything, studentID, -1).Return(nil)
		mockRepo.On("InsertSpend", mock.Anything, mock.Anything, studentID, essayID, 1).Return(created, nil)

		app, err := setupTestAppWithEssayReview(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodPost, "/reviews", map[string]any{
			"student_id": studentID, "essay_id": essayID, "amount": 1,
		})

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var out models.EssayReviewTransaction
		require.NoError(t, json.Unmarshal(body, &out))
		assert.Equal(t, -1, out.Subtotal)
		assert.Equal(t, models.ReviewEntrySpend, out.EntryType)
		require.NotNil(t, out.Status)
		assert.Equal(t, models.ReviewStatusOpen, *out.Status)
	})

	t.Run("amount below the minimum is rejected by validation", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayReviewRepository(t)
		// The handler is never reached.

		app, err := setupTestAppWithEssayReview(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodPost, "/reviews", map[string]any{
			"student_id": uuid.New(), "essay_id": uuid.New(), "amount": 0,
		})

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Contains(t, string(body), "expected number >= 1")
	})

	t.Run("essay with an open review returns conflict", func(t *testing.T) {
		t.Parallel()

		studentID, essayID, openID := uuid.New(), uuid.New(), uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).
			Return(&models.Student{ID: studentID, ReviewBalance: 5}, nil)
		mockRepo.On("FindOpenReviewForEssay", mock.Anything, mock.Anything, essayID).Return(&openID, nil)

		app, err := setupTestAppWithEssayReview(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodPost, "/reviews", map[string]any{
			"student_id": studentID, "essay_id": essayID, "amount": 1,
		})

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		assert.Contains(t, string(body), "essay already has an open review")
	})
}

func TestRoute_CompleteEssayReview(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		id, studentID, essayID := uuid.New(), uuid.New(), uuid.New()
		openStatus := models.ReviewStatusOpen
		completedStatus := models.ReviewStatusCompleted
		completedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockTransaction", mock.Anything, mock.Anything, id).Return(
			&models.EssayReviewTransaction{ID: id, StudentID: studentID, EssayID: &essayID,
				EntryType: models.ReviewEntrySpend, Status: &openStatus}, nil)
		mockRepo.On("MarkCompleted", mock.Anything, mock.Anything, id).Return(
			&models.EssayReviewTransaction{ID: id, StudentID: studentID, EssayID: &essayID,
				EntryType: models.ReviewEntrySpend, Status: &completedStatus, CompletedAt: &completedAt}, nil)

		app, err := setupTestAppWithEssayReview(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodPost, "/reviews/"+id.String()+"/complete", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var out models.EssayReviewTransaction
		require.NoError(t, json.Unmarshal(body, &out))
		require.NotNil(t, out.Status)
		assert.Equal(t, models.ReviewStatusCompleted, *out.Status)
		assert.NotNil(t, out.CompletedAt)
	})

	t.Run("already completed returns conflict", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		completed := models.ReviewStatusCompleted

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockTransaction", mock.Anything, mock.Anything, id).Return(
			&models.EssayReviewTransaction{ID: id, EntryType: models.ReviewEntrySpend, Status: &completed}, nil)

		app, err := setupTestAppWithEssayReview(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodPost, "/reviews/"+id.String()+"/complete", nil)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		assert.Contains(t, string(body), "already been completed")
	})

	t.Run("malformed id is rejected by validation", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewEssayReviewRepository(t)

		app, err := setupTestAppWithEssayReview(mockRepo)
		require.NoError(t, err)

		resp, _ := doJSON(t, app, http.MethodPost, "/reviews/not-a-uuid/complete", nil)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRoute_RefundEssayReview(t *testing.T) {
	t.Parallel()

	t.Run("success returns the charge and its reversal", func(t *testing.T) {
		t.Parallel()

		id, studentID, essayID, reversalID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
		openStatus := models.ReviewStatusOpen
		refundedStatus := models.ReviewStatusRefunded

		charge := &models.EssayReviewTransaction{ID: id, StudentID: studentID, EssayID: &essayID,
			Subtotal: -1, EntryType: models.ReviewEntrySpend, Status: &openStatus}
		reversal := &models.EssayReviewTransaction{ID: reversalID, StudentID: studentID, EssayID: &essayID,
			Subtotal: 1, EntryType: models.ReviewEntryRefund}
		refunded := &models.EssayReviewTransaction{ID: id, StudentID: studentID, EssayID: &essayID,
			Subtotal: -1, EntryType: models.ReviewEntrySpend, Status: &refundedStatus, Refund: &reversalID}

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockTransaction", mock.Anything, mock.Anything, id).Return(charge, nil)
		mockRepo.On("LockStudent", mock.Anything, mock.Anything, studentID).Return(&models.Student{ID: studentID}, nil)
		mockRepo.On("InsertRefund", mock.Anything, mock.Anything, charge).Return(reversal, nil)
		mockRepo.On("LinkRefund", mock.Anything, mock.Anything, id, reversalID).Return(refunded, nil)
		mockRepo.On("AdjustStudentBalance", mock.Anything, mock.Anything, studentID, 1).Return(nil)

		app, err := setupTestAppWithEssayReview(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodPost, "/reviews/"+id.String()+"/refund", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var out models.RefundEssayReviewResponseBody
		require.NoError(t, json.Unmarshal(body, &out))
		require.NotNil(t, out.Transaction.Status)
		assert.Equal(t, models.ReviewStatusRefunded, *out.Transaction.Status)
		assert.Equal(t, 1, out.Refund.Subtotal)
		assert.Equal(t, models.ReviewEntryRefund, out.Refund.EntryType)
	})

	t.Run("a reversal row cannot itself be refunded", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		expectTx(mockRepo)
		mockRepo.On("LockTransaction", mock.Anything, mock.Anything, id).Return(
			&models.EssayReviewTransaction{ID: id, Subtotal: 1, EntryType: models.ReviewEntryRefund}, nil)

		app, err := setupTestAppWithEssayReview(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodPost, "/reviews/"+id.String()+"/refund", nil)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, string(body), "entry_type is refund")
	})
}

func TestRoute_GetEssayReviewStatus(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		essayID, transactionID, studentID := uuid.New(), uuid.New(), uuid.New()
		expected := &models.EssayReviewStatus{
			TransactionID: transactionID, EssayID: essayID, StudentID: studentID,
			Status: models.ReviewStatusOpen, RequestedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		mockRepo := mocks.NewEssayReviewRepository(t)
		mockRepo.On("DB").Return(nil)
		mockRepo.On("FindEssayReviewStatus", mock.Anything, mock.Anything, essayID).Return(expected, nil)

		app, err := setupTestAppWithEssayReview(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/essays/"+essayID.String()+"/review-status", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var out models.EssayReviewStatus
		require.NoError(t, json.Unmarshal(body, &out))
		assert.Equal(t, transactionID, out.TransactionID)
		assert.Equal(t, models.ReviewStatusOpen, out.Status)
	})

	t.Run("unreviewed essay returns not found", func(t *testing.T) {
		t.Parallel()

		essayID := uuid.New()

		mockRepo := mocks.NewEssayReviewRepository(t)
		mockRepo.On("DB").Return(nil)
		mockRepo.On("FindEssayReviewStatus", mock.Anything, mock.Anything, essayID).
			Return(nil, errs.NotFound("essay review", "essay_id", essayID.String()))

		app, err := setupTestAppWithEssayReview(mockRepo)
		require.NoError(t, err)

		resp, _ := doJSON(t, app, http.MethodGet, "/essays/"+essayID.String()+"/review-status", nil)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
