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
	"inspirate-consulting/internal/pagination"
)

func setupTestAppWithStudent(mockRepo data.EssayReviewRepository) (*fiber.App, error) {
	cfg := config.Config{TestMode: true}
	repo := &data.Repository{EssayReview: mockRepo}
	app, _, err := SetupApp(cfg, repo, nil)
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

func setupTestAppWithStudentRepo(mockRepo data.StudentRepository) (*fiber.App, error) {
	cfg := config.Config{TestMode: true}
	repo := &data.Repository{Student: mockRepo}
	app, _, err := SetupApp(cfg, repo, nil)
	return app, err
}

func TestRoute_AssignCounselor(t *testing.T) {
	t.Parallel()

	t.Run("success returns the student with the new counselor", func(t *testing.T) {
		t.Parallel()

		studentUserID := uuid.New()
		counselorID := uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("AssignCounselor", mock.Anything, studentUserID, counselorID).
			Return(&models.Student{ID: uuid.New(), UserID: studentUserID, Year: models.Junior, GPA: 4, CounselorID: &counselorID}, nil)

		app, err := setupTestAppWithStudentRepo(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodPatch, "/students/assign",
			map[string]any{"user_id": studentUserID, "counselor_id": counselorID})

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var out models.Student
		require.NoError(t, json.Unmarshal(body, &out))
		assert.Equal(t, studentUserID, out.UserID)
		require.NotNil(t, out.CounselorID)
		assert.Equal(t, counselorID, *out.CounselorID)
	})

	t.Run("unknown student returns not found", func(t *testing.T) {
		t.Parallel()

		studentUserID := uuid.New()
		counselorID := uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("AssignCounselor", mock.Anything, studentUserID, counselorID).
			Return(nil, errs.NotFound("student", "user_id", studentUserID))

		app, err := setupTestAppWithStudentRepo(mockRepo)
		require.NoError(t, err)

		resp, _ := doJSON(t, app, http.MethodPatch, "/students/assign",
			map[string]any{"user_id": studentUserID, "counselor_id": counselorID})

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("unknown counselor returns not found", func(t *testing.T) {
		t.Parallel()

		studentUserID := uuid.New()
		counselorID := uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("AssignCounselor", mock.Anything, studentUserID, counselorID).
			Return(nil, errs.NotFound("counselor", "id", counselorID))

		app, err := setupTestAppWithStudentRepo(mockRepo)
		require.NoError(t, err)

		resp, _ := doJSON(t, app, http.MethodPatch, "/students/assign",
			map[string]any{"user_id": studentUserID, "counselor_id": counselorID})

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("assigning a counselor to another counselor returns not found", func(t *testing.T) {
		t.Parallel()

		// A counselor's user id has no student row, so there is no student to assign.
		counselorUserID := uuid.New()
		counselorID := uuid.New()

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("AssignCounselor", mock.Anything, counselorUserID, counselorID).
			Return(nil, errs.NotFound("student", "user_id", counselorUserID))

		app, err := setupTestAppWithStudentRepo(mockRepo)
		require.NoError(t, err)

		resp, _ := doJSON(t, app, http.MethodPatch, "/students/assign",
			map[string]any{"user_id": counselorUserID, "counselor_id": counselorID})

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("malformed ids are rejected by validation", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewStudentRepository(t)
		// The handler is never reached.

		app, err := setupTestAppWithStudentRepo(mockRepo)
		require.NoError(t, err)

		resp, _ := doJSON(t, app, http.MethodPatch, "/students/assign",
			map[string]any{"user_id": "not-a-uuid", "counselor_id": "also-not-a-uuid"})

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func listNameIs(want string) any {
	return mock.MatchedBy(func(name *string) bool { return name != nil && *name == want })
}

func listNoCursor() any {
	return mock.MatchedBy(func(after *pagination.NameIDKey) bool { return after == nil })
}

func listCursorAt(name string, id uuid.UUID) any {
	return mock.MatchedBy(func(after *pagination.NameIDKey) bool {
		return after != nil && after.Name == name && after.ID == id
	})
}

func studentListRows(names ...string) []models.StudentListItem {
	rows := make([]models.StudentListItem, len(names))
	for i, name := range names {
		rows[i] = models.StudentListItem{ID: uuid.New(), UserID: uuid.New(), Name: name, Year: models.Junior}
	}
	return rows
}

func decodeStudentPage(t *testing.T, body []byte) ([]string, *string) {
	t.Helper()

	var out models.FetchStudentsBody
	require.NoError(t, json.Unmarshal(body, &out))
	names := make([]string, len(out.Students))
	for i, s := range out.Students {
		names[i] = s.Name
	}
	return names, out.NextCursor
}

func TestRoute_FetchStudents(t *testing.T) {
	t.Parallel()

	t.Run("no query params uses the default limit of 5", func(t *testing.T) {
		t.Parallel()

		rows := studentListRows("Anna Lee", "Ben Cho", "Cara Diaz", "Dan Wu", "Eva Park", "Finn Ross")

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("ListStudents", mock.Anything, listNameIs(""), listNoCursor(), 6).Return(rows, nil)

		app, err := setupTestAppWithStudentRepo(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/students/fetchall", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, next := decodeStudentPage(t, body)
		assert.Equal(t, []string{"Anna Lee", "Ben Cho", "Cara Diaz", "Dan Wu", "Eva Park"}, names)
		assert.NotNil(t, next)
	})

	t.Run("a custom limit is passed through", func(t *testing.T) {
		t.Parallel()

		rows := studentListRows("Anna Lee", "Ben Cho", "Cara Diaz", "Dan Wu")

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("ListStudents", mock.Anything, listNameIs(""), listNoCursor(), 4).Return(rows, nil)

		app, err := setupTestAppWithStudentRepo(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/students/fetchall?limit=3", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, next := decodeStudentPage(t, body)
		assert.Equal(t, []string{"Anna Lee", "Ben Cho", "Cara Diaz"}, names)
		assert.NotNil(t, next)
	})

	t.Run("a name is passed through as the search term", func(t *testing.T) {
		t.Parallel()

		rows := studentListRows("Anna Lee", "Hannah Kim")

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("ListStudents", mock.Anything, listNameIs("ann"), listNoCursor(), 6).Return(rows, nil)

		app, err := setupTestAppWithStudentRepo(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/students/fetchall?name=ann", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, next := decodeStudentPage(t, body)
		assert.Equal(t, []string{"Anna Lee", "Hannah Kim"}, names)
		assert.Nil(t, next)
	})

	t.Run("the next cursor fetches the following page", func(t *testing.T) {
		t.Parallel()

		first := studentListRows("Anna Lee", "Ben Cho", "Cara Diaz")
		second := studentListRows("Dan Wu", "Eva Park")

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("ListStudents", mock.Anything, listNameIs(""), listNoCursor(), 3).Return(first, nil).Once()
		mockRepo.On("ListStudents", mock.Anything, listNameIs(""), listCursorAt("Ben Cho", first[1].ID), 3).Return(second, nil).Once()

		app, err := setupTestAppWithStudentRepo(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/students/fetchall?limit=2", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, next := decodeStudentPage(t, body)
		assert.Equal(t, []string{"Anna Lee", "Ben Cho"}, names)
		require.NotNil(t, next)

		resp, body = doJSON(t, app, http.MethodGet, "/students/fetchall?limit=2&cursor="+*next, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, _ = decodeStudentPage(t, body)
		assert.Equal(t, []string{"Dan Wu", "Eva Park"}, names)
	})

	t.Run("fewer rows than the limit returns a null next cursor", func(t *testing.T) {
		t.Parallel()

		rows := studentListRows("Dan Wu", "Eva Park", "Finn Ross")

		mockRepo := mocks.NewStudentRepository(t)
		mockRepo.On("ListStudents", mock.Anything, listNameIs(""), listNoCursor(), 6).Return(rows, nil)

		app, err := setupTestAppWithStudentRepo(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/students/fetchall", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, next := decodeStudentPage(t, body)
		assert.Len(t, names, 3)
		assert.Nil(t, next)
	})
}
