package routes

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// expectParent stubs the ownership lookup every note/achievement route performs.
func expectParent(m *mocks.ExtracurricularRepository, id int64) {
	m.On("GetExtracurricular", mock.Anything, id).Return(&models.ExtracurricularDetail{
		Extracurricular: models.Extracurricular{ID: id, StudentID: uuid.NewString()},
	}, nil)
}

// Route tests for /extracurriculars/{id}/notes.
func TestRoute_ExtracurricularNotes(t *testing.T) {
	t.Parallel()

	t.Run("create returns 201 with the note", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		expectParent(mockRepo, 7)
		mockRepo.On("CreateExtracurricularNote", mock.Anything, int64(7), models.ExtracurricularNoteRequest{Body: "Met with coach"}).
			Return(&models.ExtracurricularNote{ID: 1, ExtracurricularID: 7, Body: "Met with coach"}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodPost, "/extracurriculars/7/notes", map[string]any{"body": "Met with coach"})
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(body))

		var note models.ExtracurricularNote
		require.NoError(t, json.Unmarshal(body, &note))
		assert.Equal(t, int64(1), note.ID)
		assert.Equal(t, int64(7), note.ExtracurricularID)
		assert.Equal(t, "Met with coach", note.Body)
	})

	t.Run("update returns 200 with the note", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		expectParent(mockRepo, 7)
		mockRepo.On("UpdateExtracurricularNote", mock.Anything, int64(7), int64(3), models.ExtracurricularNoteRequest{Body: "Edited"}).
			Return(&models.ExtracurricularNote{ID: 3, ExtracurricularID: 7, Body: "Edited"}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodPatch, "/extracurriculars/7/notes/3", map[string]any{"body": "Edited"})
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
		assert.Contains(t, string(body), `"body":"Edited"`)
	})

	t.Run("delete returns 204", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		expectParent(mockRepo, 7)
		mockRepo.On("DeleteExtracurricularNote", mock.Anything, int64(7), int64(3)).Return(nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodDelete, "/extracurriculars/7/notes/3", nil)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode, string(body))
		assert.Empty(t, body)
	})

	t.Run("note not found under this extracurricular is 404", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		expectParent(mockRepo, 7)
		mockRepo.On("UpdateExtracurricularNote", mock.Anything, int64(7), int64(99), mock.Anything).Return(nil, errs.NotFound("note", "id", int64(99)))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPatch, "/extracurriculars/7/notes/99", map[string]any{"body": "x"})
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("missing parent is 404 and the note is not touched", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(404)).Return(nil, errs.NotFound("extracurricular", "id", int64(404)))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/extracurriculars/404/notes", map[string]any{"body": "x"})
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	validationCases := []struct {
		name, method, path string
		body               map[string]any
	}{
		{"create with empty body text", http.MethodPost, "/extracurriculars/7/notes", map[string]any{"body": ""}},
		{"create without body field", http.MethodPost, "/extracurriculars/7/notes", map[string]any{}},
		{"create with wrong type", http.MethodPost, "/extracurriculars/7/notes", map[string]any{"body": 5}},
		{"update with empty body text", http.MethodPatch, "/extracurriculars/7/notes/3", map[string]any{"body": ""}},
		{"non-numeric extracurricular id", http.MethodPost, "/extracurriculars/abc/notes", map[string]any{"body": "x"}},
		{"non-numeric note id", http.MethodPatch, "/extracurriculars/7/notes/abc", map[string]any{"body": "x"}},
	}
	for _, tc := range validationCases {
		t.Run("validation error - "+tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewExtracurricularRepository(t)
			app, err := setupExtracurricularTestApp(mockRepo)
			require.NoError(t, err)

			resp, body := doRequest(t, app, tc.method, tc.path, tc.body)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, string(body))
		})
	}
}

// Route tests for /extracurriculars/{id}/achievements.
func TestRoute_ExtracurricularAchievements(t *testing.T) {
	t.Parallel()

	t.Run("create returns 201 with the achievement", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		expectParent(mockRepo, 7)
		mockRepo.On("CreateExtracurricularAchievement", mock.Anything, int64(7),
			models.CreateExtracurricularAchievementRequest{EntryDate: "2025-03-15", Body: "Won regionals"}).
			Return(&models.ExtracurricularAchievement{
				ID: 2, ExtracurricularID: 7, Body: "Won regionals",
				EntryDate: time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
			}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodPost, "/extracurriculars/7/achievements",
			map[string]any{"entry_date": "2025-03-15", "body": "Won regionals"})
		require.Equal(t, http.StatusCreated, resp.StatusCode, string(body))

		var a models.ExtracurricularAchievement
		require.NoError(t, json.Unmarshal(body, &a))
		assert.Equal(t, int64(2), a.ID)
		assert.Equal(t, "Won regionals", a.Body)
	})

	t.Run("partial update forwards only sent fields", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		expectParent(mockRepo, 7)
		mockRepo.On("UpdateExtracurricularAchievement", mock.Anything, int64(7), int64(4), mock.MatchedBy(func(b models.UpdateExtracurricularAchievementRequest) bool {
			return b.EntryDate == nil && b.Body != nil && *b.Body == "Edited"
		})).Return(&models.ExtracurricularAchievement{ID: 4, ExtracurricularID: 7, Body: "Edited"}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodPatch, "/extracurriculars/7/achievements/4", map[string]any{"body": "Edited"})
		assert.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	})

	t.Run("delete returns 204", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		expectParent(mockRepo, 7)
		mockRepo.On("DeleteExtracurricularAchievement", mock.Anything, int64(7), int64(4)).Return(nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodDelete, "/extracurriculars/7/achievements/4", nil)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode, string(body))
	})

	t.Run("delete of missing achievement is 404", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		expectParent(mockRepo, 7)
		mockRepo.On("DeleteExtracurricularAchievement", mock.Anything, int64(7), int64(99)).Return(errs.NotFound("achievement", "id", int64(99)))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/extracurriculars/7/achievements/99", nil)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	validationCases := []struct {
		name, method, path string
		body               map[string]any
	}{
		{"create without entry_date", http.MethodPost, "/extracurriculars/7/achievements", map[string]any{"body": "x"}},
		{"create with malformed entry_date", http.MethodPost, "/extracurriculars/7/achievements", map[string]any{"entry_date": "03/15/2025", "body": "x"}},
		{"create with empty body text", http.MethodPost, "/extracurriculars/7/achievements", map[string]any{"entry_date": "2025-03-15", "body": ""}},
		{"update with malformed entry_date", http.MethodPatch, "/extracurriculars/7/achievements/4", map[string]any{"entry_date": "soon"}},
		{"update with empty body text", http.MethodPatch, "/extracurriculars/7/achievements/4", map[string]any{"body": ""}},
		{"non-numeric achievement id", http.MethodDelete, "/extracurriculars/7/achievements/abc", nil},
	}
	for _, tc := range validationCases {
		t.Run("validation error - "+tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewExtracurricularRepository(t)
			app, err := setupExtracurricularTestApp(mockRepo)
			require.NoError(t, err)

			var body any
			if tc.body != nil {
				body = tc.body
			}
			resp, respBody := doRequest(t, app, tc.method, tc.path, body)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, string(respBody))
		})
	}
}
