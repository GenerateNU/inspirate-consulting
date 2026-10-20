// test hitting get, create user endpoints
package routes

import (
	"bytes"
	"encoding/json"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/gofiber/fiber/v3"
)

func setupTestAppWithUser(mockUser data.UserRepository, supabase supabase.SupabaseInterface) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
		Supabase: supabase,
	}
	repo := &data.Repository{
		User: mockUser,
	}
	app, _, err := SetupApp(cfg, repo, nil)
	return app, err
}

func TestRoute_CreateUser(t *testing.T) {
	// creates user

	t.Parallel()
	u1 := uuid.New()
	s1 := uuid.New()

	key := "key"
	email := "alengano123@gmail.com"
	name := "Aleng123"

	userTest := models.User{
		ID:         u1,
		Name:       "Aleng123",
		SupabaseID: s1,
		PfpKey:     &key,
	}

	t.Run("creation success", func(t *testing.T) {

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(in models.CreateUserInput) bool {
			return in.Body.Name == "Aleng123"
		}), mock.Anything).Return(&models.CreateUserOutput{
			Body: &models.CreateUserBody{User: &userTest},
		}, nil)

		app, err := setupTestAppWithUser(mockRepo, &supabase.MockSupabase{})
		println(app)
		require.NoError(t, err)

		// Send valid JSON payload
		payload := map[string]any{
			"email":   email,
			"name":    name,
			"pfp_key": &key,
			"role":    "student",
		}
		bodyBytes, err := json.Marshal(payload)
		println(bodyBytes)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		// Verify 200 OK status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify returned JSON body
		respBody, err := io.ReadAll(resp.Body)
		println(respBody)
		require.NoError(t, err)

		var output models.CreateUserOutput
		err = json.Unmarshal(respBody, &output.Body)
		require.NoError(t, err)

		t.Log("value:", output)
		assert.Equal(t, name, output.Body.User.Name)

	})

	t.Run("failure (sending empty string to name)", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewUserRepository(t)

		app, err := setupTestAppWithUser(mockRepo, &supabase.MockSupabase{})
		require.NoError(t, err)

		payload := map[string]any{
			"name": "",
		}

		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	})
}

func TestRoute_FetchUser(t *testing.T) {
	t.Parallel()
	test_uuid := uuid.MustParse("7bebfe6e-36ca-4343-a46c-d9430a951ba7")
	//test_uuid_two := uuid.MustParse("ba0ccc36-6c60-4ad3-90ad-0fed4378f3df")
	s1 := uuid.New()
	key := ""

	userTest := models.User{
		ID:         test_uuid,
		Name:       "Aleng123",
		SupabaseID: s1,
		PfpKey:     &key,
	}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewUserRepository(t)
		input := models.FetchUserInput{ID: test_uuid}
		mockRepo.On("FetchUser", mock.Anything, input).Return(&models.FetchUserOutput{Body: &userTest}, nil)
		app, err := setupTestAppWithUser(mockRepo, &supabase.MockSupabase{})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/user/"+test_uuid.String(), nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.FetchUserOutput
		err = json.Unmarshal(respBody, &output.Body)
		require.NoError(t, err)
		assert.Equal(t, test_uuid, output.Body.ID)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewUserRepository(t)

		input := models.FetchUserInput{ID: test_uuid}

		mockRepo.On("FetchUser", mock.Anything, input).
			Return(nil, errs.NotFound("user", "user id", test_uuid.String()))

		app, err := setupTestAppWithUser(mockRepo, &supabase.MockSupabase{})
		println(app)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/user/"+test_uuid.String(), nil)
		require.NoError(t, err)

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		// User repository returns an errs.NotFound (HTTPError).
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestRoute_FetchCurrentUser(t *testing.T) {
	t.Parallel()

	t.Run("returns the user with their counselor id", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()
		counselorID := uuid.New()

		mockUser := mocks.NewUserRepository(t)
		mockUser.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
			Return(&models.FetchUserOutput{Body: &models.User{ID: userID, Name: "Test Counselor"}}, nil)

		app, err := setupTestAppWithUser(mockUser, &supabase.MockSupabase{})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/user/me", nil)
		require.NoError(t, err)
		req.Header.Set("X-User-ID", userID.String())
		req.Header.Set("X-Counselor-ID", counselorID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var out models.SessionUser
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
		assert.Equal(t, userID, out.ID)
		assert.Equal(t, "Test Counselor", out.Name)
		require.NotNil(t, out.CounselorID)
		assert.Equal(t, counselorID, *out.CounselorID)
		assert.Nil(t, out.StudentID)
	})

	t.Run("unknown user returns not found", func(t *testing.T) {
		t.Parallel()

		userID := uuid.New()

		mockUser := mocks.NewUserRepository(t)
		mockUser.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
			Return(nil, errs.NotFound("user", "id", userID.String()))

		app, err := setupTestAppWithUser(mockUser, &supabase.MockSupabase{})
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodGet, "/user/me", nil)
		require.NoError(t, err)
		req.Header.Set("X-User-ID", userID.String())

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
