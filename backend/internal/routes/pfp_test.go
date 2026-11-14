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

// creates a Fiber app in test mode with mocked user and pfp repositories.
func setupTestAppWithPfp(t *testing.T, mockUser data.UserRepository, mockPfp data.PfpRepository) *fiber.App {
	t.Helper()
	repo := &data.Repository{User: mockUser, Pfp: mockPfp}
	app, _, err := SetupApp(config.Config{TestMode: true}, repo, nil, nil)
	require.NoError(t, err)
	return app
}

func TestRoute_PresignPfpUpload(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	asUser := map[string]string{"X-User-ID": userID.String()}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		key := "pfps/" + userID.String() + "/a.jpg"
		mockPfp := mocks.NewPfpRepository(t)
		mockPfp.On("PresignUpload", mock.Anything, userID, models.PresignPfpUploadRequestBody{
			ContentType:   "image/jpeg",
			ContentLength: 2048,
		}).Return(&models.PresignUploadResponse{S3Key: key, UploadURL: "https://bucket/upload"}, nil)

		app := setupTestAppWithPfp(t, mocks.NewUserRepository(t), mockPfp)
		resp, body := doJSON(t, app, http.MethodPut, "/user/pfp", map[string]any{
			"content_type":   "image/jpeg",
			"content_length": 2048,
		}, asUser)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var out models.PresignUploadResponse
		require.NoError(t, json.Unmarshal(body, &out))
		assert.Equal(t, key, out.S3Key)
	})

	t.Run("validation errors", func(t *testing.T) {
		t.Parallel()

		cases := map[string]map[string]any{
			"missing content_type":    {"content_length": 2048},
			"non-image content_type":  {"content_type": "video/mp4", "content_length": 2048},
			"svg content_type":        {"content_type": "image/svg+xml", "content_length": 2048},
			"missing content_length":  {"content_type": "image/png"},
			"zero content_length":     {"content_type": "image/png", "content_length": 0},
			"content_length over 5MB": {"content_type": "image/png", "content_length": 5242881},
		}
		for name, payload := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				// Mock should not be called when Huma validation fails.
				app := setupTestAppWithPfp(t, mocks.NewUserRepository(t), mocks.NewPfpRepository(t))
				resp, _ := doJSON(t, app, http.MethodPut, "/user/pfp", payload, asUser)

				assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
			})
		}
	})
}

func TestRoute_UpdatePfp(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	asUser := map[string]string{"X-User-ID": userID.String()}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		key := "pfps/" + userID.String() + "/a.jpg"
		mockPfp := mocks.NewPfpRepository(t)
		mockPfp.On("VerifyUpload", mock.Anything, userID, key).Return(nil)
		mockPfp.On("GetDownloadURL", mock.Anything, key).Return("https://bucket/a.jpg", nil)
		mockUser := mocks.NewUserRepository(t)
		mockUser.On("SetPfpKey", mock.Anything, userID, &key).Return(&models.User{ID: userID, PfpKey: &key}, nil, nil)

		app := setupTestAppWithPfp(t, mockUser, mockPfp)
		resp, body := doJSON(t, app, http.MethodPatch, "/user/pfp", map[string]any{"s3_key": key}, asUser)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var out models.Pfp
		require.NoError(t, json.Unmarshal(body, &out))
		assert.Equal(t, key, out.S3Key)
		assert.Equal(t, "https://bucket/a.jpg", out.DownloadURL)
	})

	t.Run("another user's key is forbidden", func(t *testing.T) {
		t.Parallel()

		otherKey := "pfps/" + uuid.NewString() + "/theirs.jpg"
		mockPfp := mocks.NewPfpRepository(t)
		mockPfp.On("VerifyUpload", mock.Anything, userID, otherKey).
			Return(errs.Forbidden("you can only use your own profile picture"))

		app := setupTestAppWithPfp(t, mocks.NewUserRepository(t), mockPfp)
		resp, _ := doJSON(t, app, http.MethodPatch, "/user/pfp", map[string]any{"s3_key": otherKey}, asUser)

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("validation error : missing s3_key", func(t *testing.T) {
		t.Parallel()

		app := setupTestAppWithPfp(t, mocks.NewUserRepository(t), mocks.NewPfpRepository(t))
		resp, _ := doJSON(t, app, http.MethodPatch, "/user/pfp", map[string]any{}, asUser)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestRoute_DeletePfp(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	asUser := map[string]string{"X-User-ID": userID.String()}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		key := "pfps/" + userID.String() + "/a.jpg"
		mockUser := mocks.NewUserRepository(t)
		mockUser.On("SetPfpKey", mock.Anything, userID, (*string)(nil)).Return(&models.User{ID: userID}, &key, nil)
		mockPfp := mocks.NewPfpRepository(t)
		mockPfp.On("DeletePfp", mock.Anything, userID, key).Return(nil)

		app := setupTestAppWithPfp(t, mockUser, mockPfp)
		resp, _ := doJSON(t, app, http.MethodDelete, "/user/pfp", nil, asUser)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("no picture returns not found", func(t *testing.T) {
		t.Parallel()

		mockUser := mocks.NewUserRepository(t)
		mockUser.On("SetPfpKey", mock.Anything, userID, (*string)(nil)).Return(&models.User{ID: userID}, nil, nil)

		app := setupTestAppWithPfp(t, mockUser, mocks.NewPfpRepository(t))
		resp, _ := doJSON(t, app, http.MethodDelete, "/user/pfp", nil, asUser)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
