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
	"inspirate-consulting/internal/models"
)

func setupTestAppWithCounselorRepo(mockRepo data.CounselorRepository) (*fiber.App, error) {
	cfg := config.Config{TestMode: true}
	repo := &data.Repository{Counselor: mockRepo}
	app, _, err := SetupApp(cfg, repo, nil)
	return app, err
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

		app, err := setupTestAppWithCounselorRepo(mockRepo)
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

		app, err := setupTestAppWithCounselorRepo(mockRepo)
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

		app, err := setupTestAppWithCounselorRepo(mockRepo)
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

		app, err := setupTestAppWithCounselorRepo(mockRepo)
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

		app, err := setupTestAppWithCounselorRepo(mockRepo)
		require.NoError(t, err)

		resp, body := doJSON(t, app, http.MethodGet, "/counselors/fetchall", nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		names, next := decodeCounselorPage(t, body)
		assert.Len(t, names, 2)
		assert.Nil(t, next)
	})
}
