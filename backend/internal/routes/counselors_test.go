package routes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

// setupTestAppWithCounselor creates a Fiber app in test mode with a mocked counselor repository
func setupTestAppWithCounselor(mockCounselorRepo data.CounselorRepository) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	repo := &data.Repository{
		Counselor: mockCounselorRepo,
	}
	app, _, err := SetupApp(cfg, repo, nil)
	return app, err
}

func TestRoute_ListCounselorPendingReviews(t *testing.T) {
	t.Parallel()

	counselorID := uuid.New()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		reviewID := uuid.New()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListPendingReviews", mock.Anything, counselorID).
			Return([]models.PendingReview{
				{TransactionID: reviewID, EssayID: uuid.New(), EssayType: "personal-statement", EssayLink: "https://example.com/a", StudentID: uuid.New(), StudentName: "Alice"},
			}, nil)

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/reviews", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.PendingReview
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.Len(t, output, 1)
		assert.Equal(t, reviewID, output[0].TransactionID)
		assert.Equal(t, "Alice", output[0].StudentName)
		assert.Equal(t, "https://example.com/a", output[0].EssayLink)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListPendingReviews", mock.Anything, counselorID).
			Return([]models.PendingReview{}, nil)

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/reviews", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.PendingReview
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		assert.Len(t, output, 0)
	})

	t.Run("repository error returns internal server error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListPendingReviews", mock.Anything, counselorID).
			Return(nil, errors.New("database error"))

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/reviews", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestRoute_ListCounselorStudents(t *testing.T) {
	t.Parallel()

	counselorID := uuid.New()

	t.Run("success applies default sort and sets flags", func(t *testing.T) {
		t.Parallel()

		// Alice has a pending review, Bob is falling behind, Carol neither
		aliceID, bobID, carolID := uuid.New(), uuid.New(), uuid.New()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListStudents", mock.Anything, counselorID).
			Return([]models.CounselorStudentRow{
				{StudentID: carolID, Name: "Carol", Year: models.Freshman},
				{StudentID: bobID, Name: "Bob", Year: models.Senior},
				{StudentID: aliceID, Name: "Alice", Year: models.Junior},
			}, nil)
		mockRepo.On("ListPendingReviews", mock.Anything, counselorID).
			Return([]models.PendingReview{{TransactionID: uuid.New(), StudentID: aliceID, StudentName: "Alice"}}, nil)
		mockRepo.On("ListFallingBehindStudentIDs", mock.Anything, counselorID).
			Return([]uuid.UUID{bobID}, nil)

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/students", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.CounselorStudent
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.Len(t, output, 3)

		// sort_by and sort_order default to name asc
		assert.Equal(t, "Alice", output[0].Name)
		assert.Equal(t, aliceID, output[0].StudentID)
		assert.True(t, output[0].HasPendingReview)
		assert.False(t, output[0].IsFallingBehind)

		assert.Equal(t, "Bob", output[1].Name)
		assert.False(t, output[1].HasPendingReview)
		assert.True(t, output[1].IsFallingBehind)

		assert.Equal(t, "Carol", output[2].Name)
		assert.False(t, output[2].HasPendingReview)
		assert.False(t, output[2].IsFallingBehind)
	})

	t.Run("filters by pending review", func(t *testing.T) {
		t.Parallel()

		aliceID, bobID := uuid.New(), uuid.New()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListStudents", mock.Anything, counselorID).
			Return([]models.CounselorStudentRow{
				{StudentID: bobID, Name: "Bob", Year: models.Senior},
				{StudentID: aliceID, Name: "Alice", Year: models.Junior},
			}, nil)
		mockRepo.On("ListPendingReviews", mock.Anything, counselorID).
			Return([]models.PendingReview{{TransactionID: uuid.New(), StudentID: aliceID, StudentName: "Alice"}}, nil)
		mockRepo.On("ListFallingBehindStudentIDs", mock.Anything, counselorID).
			Return([]uuid.UUID{}, nil)

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/students?has_pending_review=true", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.CounselorStudent
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.Len(t, output, 1)
		assert.Equal(t, "Alice", output[0].Name)
	})

	t.Run("filters by year", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListStudents", mock.Anything, counselorID).
			Return([]models.CounselorStudentRow{
				{StudentID: uuid.New(), Name: "Bob", Year: models.Senior},
				{StudentID: uuid.New(), Name: "Alice", Year: models.Junior},
			}, nil)
		mockRepo.On("ListPendingReviews", mock.Anything, counselorID).
			Return([]models.PendingReview{}, nil)
		mockRepo.On("ListFallingBehindStudentIDs", mock.Anything, counselorID).
			Return([]uuid.UUID{}, nil)

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/students?year=senior", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.CounselorStudent
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.Len(t, output, 1)
		assert.Equal(t, "Bob", output[0].Name)
	})

	t.Run("sorts by year desc", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListStudents", mock.Anything, counselorID).
			Return([]models.CounselorStudentRow{
				{StudentID: uuid.New(), Name: "Carol", Year: models.Freshman},
				{StudentID: uuid.New(), Name: "Bob", Year: models.Senior},
				{StudentID: uuid.New(), Name: "Alice", Year: models.Junior},
			}, nil)
		mockRepo.On("ListPendingReviews", mock.Anything, counselorID).
			Return([]models.PendingReview{}, nil)
		mockRepo.On("ListFallingBehindStudentIDs", mock.Anything, counselorID).
			Return([]uuid.UUID{}, nil)

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/students?sort_by=year&sort_order=desc", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output []models.CounselorStudent
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.Len(t, output, 3)
		assert.Equal(t, "Bob", output[0].Name)
		assert.Equal(t, "Alice", output[1].Name)
		assert.Equal(t, "Carol", output[2].Name)
	})

	t.Run("no matches returns an empty list", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListStudents", mock.Anything, counselorID).
			Return([]models.CounselorStudentRow{
				{StudentID: uuid.New(), Name: "Alice", Year: models.Junior},
			}, nil)
		mockRepo.On("ListPendingReviews", mock.Anything, counselorID).
			Return([]models.PendingReview{}, nil)
		mockRepo.On("ListFallingBehindStudentIDs", mock.Anything, counselorID).
			Return([]uuid.UUID{}, nil)

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/students?year=sophomore", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.JSONEq(t, "[]", string(respBody))
	})

	t.Run("validation error : invalid year", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		// Huma will reject the request before reaching the handler

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/students?year=graduate", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("validation error : invalid sort_by", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		// Huma will reject the request before reaching the handler

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/students?sort_by=gpa", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("validation error : invalid sort_order", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		// Huma will reject the request before reaching the handler

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/students?sort_order=sideways", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("validation error : invalid has_pending_review", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		// Huma will reject the request before reaching the handler

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/students?has_pending_review=maybe", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("repository error returns internal server error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListStudents", mock.Anything, counselorID).
			Return(nil, errors.New("database error"))

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/counselors/students", nil)
		require.NoError(t, err)
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func counselorListRows(names ...string) []models.CounselorListItem {
	rows := make([]models.CounselorListItem, len(names))
	for i, name := range names {
		rows[i] = models.CounselorListItem{ID: uuid.New(), UserID: uuid.New(), Name: name}
	}
	return rows
}

func decodeCounselorPage(t *testing.T, body []byte) ([]string, *string) {
	t.Helper()

	var out models.FetchCounselorsBody
	require.NoError(t, json.Unmarshal(body, &out))
	names := make([]string, len(out.Counselors))
	for i, c := range out.Counselors {
		names[i] = c.Name
	}
	return names, out.NextCursor
}

func TestRoute_FetchCounselors(t *testing.T) {
	t.Parallel()

	t.Run("no query params uses the default limit of 5", func(t *testing.T) {
		t.Parallel()

		rows := counselorListRows("David Chen", "James Okafor", "Maria Lopez", "Priya Patel", "Sam Rivera", "Tara Singh")

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListCounselors", mock.Anything, listNameIs(""), listNoCursor(), 6).Return(rows, nil)

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/counselors/fetchall", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, next := decodeCounselorPage(t, body)
		assert.Equal(t, []string{"David Chen", "James Okafor", "Maria Lopez", "Priya Patel", "Sam Rivera"}, names)
		assert.NotNil(t, next)
	})

	t.Run("a custom limit is passed through", func(t *testing.T) {
		t.Parallel()

		rows := counselorListRows("David Chen", "James Okafor", "Maria Lopez")

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListCounselors", mock.Anything, listNameIs(""), listNoCursor(), 3).Return(rows, nil)

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/counselors/fetchall?limit=2", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, next := decodeCounselorPage(t, body)
		assert.Equal(t, []string{"David Chen", "James Okafor"}, names)
		assert.NotNil(t, next)
	})

	t.Run("a name is passed through as the search term", func(t *testing.T) {
		t.Parallel()

		rows := counselorListRows("Priya Patel")

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListCounselors", mock.Anything, listNameIs("pat"), listNoCursor(), 6).Return(rows, nil)

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/counselors/fetchall?name=pat", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, next := decodeCounselorPage(t, body)
		assert.Equal(t, []string{"Priya Patel"}, names)
		assert.Nil(t, next)
	})

	t.Run("the next cursor fetches the following page", func(t *testing.T) {
		t.Parallel()

		first := counselorListRows("David Chen", "James Okafor", "Maria Lopez")
		second := counselorListRows("Maria Lopez", "Priya Patel")

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListCounselors", mock.Anything, listNameIs(""), listNoCursor(), 3).Return(first, nil).Once()
		mockRepo.On("ListCounselors", mock.Anything, listNameIs(""), listCursorAt("James Okafor", first[1].ID), 3).Return(second, nil).Once()

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/counselors/fetchall?limit=2", nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, next := decodeCounselorPage(t, body)
		assert.Equal(t, []string{"David Chen", "James Okafor"}, names)
		require.NotNil(t, next)

		resp, body = doJSON(t, app, http.MethodGet, "/counselors/fetchall?limit=2&cursor="+*next, nil)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, _ = decodeCounselorPage(t, body)
		assert.Equal(t, []string{"Maria Lopez", "Priya Patel"}, names)
	})

	t.Run("fewer rows than the limit returns a null next cursor", func(t *testing.T) {
		t.Parallel()

		rows := counselorListRows("Priya Patel", "Sam Rivera")

		mockRepo := mocks.NewCounselorRepository(t)
		mockRepo.On("ListCounselors", mock.Anything, listNameIs(""), listNoCursor(), 6).Return(rows, nil)

		app, err := setupTestAppWithCounselor(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/counselors/fetchall", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, next := decodeCounselorPage(t, body)
		assert.Len(t, names, 2)
		assert.Nil(t, next)
	})
}
