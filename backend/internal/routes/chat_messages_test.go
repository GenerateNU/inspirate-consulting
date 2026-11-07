package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// setupChatMessageTestApp creates a Fiber app in test mode with a mocked chat message repository.
func setupChatMessageTestApp(t *testing.T, mockChat data.ChatMessageRepository) *fiber.App {
	t.Helper()

	app, _, err := SetupApp(config.Config{TestMode: true}, &data.Repository{ChatMessage: mockChat}, nil, nil)
	require.NoError(t, err)
	return app
}

// Route tests for POST /messages.
func TestRoute_CreateChatMessage(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	recipientID := uuid.New()
	headers := map[string]string{"X-User-ID": userID.String()}
	payload := map[string]any{"recipient_id": recipientID, "message": "Hello"}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("CreateChatMessage", mock.Anything, userID, models.CreateChatMessageRequestBody{RecipientID: recipientID, Message: "Hello"}).
			Return(&models.ChatMessage{ID: uuid.New(), SenderID: userID, RecipientID: recipientID, Message: "Hello", CreatedAt: time.Now()}, nil)

		resp, body := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodPost, "/messages", payload, headers)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		var created models.ChatMessage
		require.NoError(t, json.Unmarshal(body, &created))
		assert.Equal(t, userID, created.SenderID)
		assert.Equal(t, recipientID, created.RecipientID)
		assert.Equal(t, "Hello", created.Message)
		assert.Nil(t, created.EditedAt)
		assert.Nil(t, created.ReadAt)
	})

	t.Run("validation errors", func(t *testing.T) {
		t.Parallel()

		cases := map[string]map[string]any{
			"missing recipient":     {"message": "Hello"},
			"missing message":       {"recipient_id": recipientID},
			"empty message":         {"recipient_id": recipientID, "message": ""},
			"non-uuid recipient id": {"recipient_id": "nope", "message": "Hello"},
		}
		for name, p := range cases {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				resp, _ := doRequest(t, setupChatMessageTestApp(t, mocks.NewChatMessageRepository(t)), http.MethodPost, "/messages", p, headers)
				assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
			})
		}
	})

	t.Run("sending to yourself returns 400", func(t *testing.T) {
		t.Parallel()

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mocks.NewChatMessageRepository(t)), http.MethodPost, "/messages",
			map[string]any{"recipient_id": userID, "message": "Hello"}, headers)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("whitespace-only message returns 400", func(t *testing.T) {
		t.Parallel()

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mocks.NewChatMessageRepository(t)), http.MethodPost, "/messages",
			map[string]any{"recipient_id": recipientID, "message": "   "}, headers)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("unknown recipient returns 404", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("CreateChatMessage", mock.Anything, userID, mock.Anything).Return(nil, errs.NotFound("user", "id", recipientID))

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodPost, "/messages", payload, headers)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("repository error returns 500", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("CreateChatMessage", mock.Anything, userID, mock.Anything).Return(nil, errors.New("database error"))

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodPost, "/messages", payload, headers)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for GET /chats.
func TestRoute_ListChats(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	headers := map[string]string{"X-User-ID": userID.String()}

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		otherID := uuid.New()
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("ListChats", mock.Anything, userID).Return([]models.ChatSummary{
			{OtherUserID: otherID, OtherUserName: "Alice", LastMessage: "Hi", LastSenderID: otherID, UnreadCount: 3},
		}, nil)

		resp, body := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodGet, "/chats", nil, headers)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		var chats []models.ChatSummary
		require.NoError(t, json.Unmarshal(body, &chats))
		require.Len(t, chats, 1)
		assert.Equal(t, otherID, chats[0].OtherUserID)
		assert.Equal(t, "Alice", chats[0].OtherUserName)
		assert.Equal(t, 3, chats[0].UnreadCount)
	})

	t.Run("repository error returns 500", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("ListChats", mock.Anything, userID).Return(nil, errors.New("database error"))

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodGet, "/chats", nil, headers)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for GET /chats/{user_id}/messages.
func TestRoute_ListChatMessages(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	otherID := uuid.New()
	headers := map[string]string{"X-User-ID": userID.String()}
	path := "/chats/" + otherID.String() + "/messages"
	nilKey := (*pagination.TimeIDKey)(nil)

	t.Run("default limit is 50", func(t *testing.T) {
		t.Parallel()

		msg := models.ChatMessage{ID: uuid.New(), SenderID: otherID, RecipientID: userID, Message: "Hi", CreatedAt: time.Now()}
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("ListChatMessages", mock.Anything, userID, otherID, nilKey, 51).Return([]models.ChatMessage{msg}, nil)

		resp, body := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodGet, path, nil, headers)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		var page models.ChatMessagePage
		require.NoError(t, json.Unmarshal(body, &page))
		require.Len(t, page.Messages, 1)
		assert.Equal(t, msg.ID, page.Messages[0].ID)
		assert.Nil(t, page.NextCursor)
	})

	t.Run("returns next cursor when more messages exist", func(t *testing.T) {
		t.Parallel()

		now := time.Now()
		rows := []models.ChatMessage{
			{ID: uuid.New(), CreatedAt: now},
			{ID: uuid.New(), CreatedAt: now.Add(-time.Minute)},
		}
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("ListChatMessages", mock.Anything, userID, otherID, nilKey, 2).Return(rows, nil)

		resp, body := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodGet, path+"?limit=1", nil, headers)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		var page models.ChatMessagePage
		require.NoError(t, json.Unmarshal(body, &page))
		assert.Len(t, page.Messages, 1)
		assert.NotNil(t, page.NextCursor)
	})

	for _, limit := range []string{"0", "101"} {
		t.Run("limit "+limit+" is rejected", func(t *testing.T) {
			t.Parallel()

			resp, _ := doRequest(t, setupChatMessageTestApp(t, mocks.NewChatMessageRepository(t)), http.MethodGet, path+"?limit="+limit, nil, headers)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		})
	}

	t.Run("invalid cursor returns 400", func(t *testing.T) {
		t.Parallel()

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mocks.NewChatMessageRepository(t)), http.MethodGet, path+"?cursor=bad!cursor", nil, headers)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("non-uuid user id is rejected", func(t *testing.T) {
		t.Parallel()

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mocks.NewChatMessageRepository(t)), http.MethodGet, "/chats/nope/messages", nil, headers)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("repository error returns 500", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("ListChatMessages", mock.Anything, userID, otherID, nilKey, 51).Return(nil, errors.New("database error"))

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodGet, path, nil, headers)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for POST /chats/{user_id}/read.
func TestRoute_MarkChatRead(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	otherID := uuid.New()
	headers := map[string]string{"X-User-ID": userID.String()}
	path := "/chats/" + otherID.String() + "/read"

	t.Run("success returns 204", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("MarkChatRead", mock.Anything, userID, otherID).Return(nil)

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodPost, path, nil, headers)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("non-uuid user id is rejected", func(t *testing.T) {
		t.Parallel()

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mocks.NewChatMessageRepository(t)), http.MethodPost, "/chats/nope/read", nil, headers)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("repository error returns 500", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("MarkChatRead", mock.Anything, userID, otherID).Return(errors.New("database error"))

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodPost, path, nil, headers)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for PATCH /messages/{id}.
func TestRoute_EditChatMessage(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	messageID := uuid.New()
	headers := map[string]string{"X-User-ID": userID.String()}
	path := "/messages/" + messageID.String()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		editedAt := time.Now()
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("EditChatMessage", mock.Anything, messageID, userID, "Edited").
			Return(&models.ChatMessage{ID: messageID, SenderID: userID, Message: "Edited", EditedAt: &editedAt}, nil)

		resp, body := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodPatch, path, map[string]any{"message": "Edited"}, headers)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		var edited models.ChatMessage
		require.NoError(t, json.Unmarshal(body, &edited))
		assert.Equal(t, "Edited", edited.Message)
		assert.NotNil(t, edited.EditedAt)
	})

	t.Run("missing message is rejected", func(t *testing.T) {
		t.Parallel()

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mocks.NewChatMessageRepository(t)), http.MethodPatch, path, map[string]any{}, headers)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("whitespace-only message returns 400", func(t *testing.T) {
		t.Parallel()

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mocks.NewChatMessageRepository(t)), http.MethodPatch, path, map[string]any{"message": "  "}, headers)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("not the sender returns 404", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("EditChatMessage", mock.Anything, messageID, userID, "Edited").
			Return(nil, errs.NotFound("chat message", "id", messageID))

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodPatch, path, map[string]any{"message": "Edited"}, headers)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("non-uuid id is rejected", func(t *testing.T) {
		t.Parallel()

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mocks.NewChatMessageRepository(t)), http.MethodPatch, "/messages/nope", map[string]any{"message": "Edited"}, headers)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

// Route tests for PATCH /messages/{id}/read.
func TestRoute_UpdateChatMessageRead(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	messageID := uuid.New()
	headers := map[string]string{"X-User-ID": userID.String()}
	path := "/messages/" + messageID.String() + "/read"

	t.Run("marking read sets read_at", func(t *testing.T) {
		t.Parallel()

		readAt := time.Now()
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("UpdateChatMessageReadAt", mock.Anything, messageID, userID, mock.MatchedBy(func(r *time.Time) bool { return r != nil })).
			Return(&models.ChatMessage{ID: messageID, RecipientID: userID, ReadAt: &readAt}, nil)

		resp, body := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodPatch, path, map[string]any{"read": true}, headers)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		var updated models.ChatMessage
		require.NoError(t, json.Unmarshal(body, &updated))
		assert.NotNil(t, updated.ReadAt)
	})

	t.Run("marking unread clears read_at", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("UpdateChatMessageReadAt", mock.Anything, messageID, userID, (*time.Time)(nil)).
			Return(&models.ChatMessage{ID: messageID, RecipientID: userID}, nil)

		resp, body := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodPatch, path, map[string]any{"read": false}, headers)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		var updated models.ChatMessage
		require.NoError(t, json.Unmarshal(body, &updated))
		assert.Nil(t, updated.ReadAt)
	})

	t.Run("missing read field is rejected", func(t *testing.T) {
		t.Parallel()

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mocks.NewChatMessageRepository(t)), http.MethodPatch, path, map[string]any{}, headers)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("not the recipient returns 404", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("UpdateChatMessageReadAt", mock.Anything, messageID, userID, mock.Anything).
			Return(nil, errs.NotFound("chat message", "id", messageID))

		resp, _ := doRequest(t, setupChatMessageTestApp(t, mockRepo), http.MethodPatch, path, map[string]any{"read": true}, headers)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
