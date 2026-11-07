package chatmessage

import (
	"context"
	"net/http"
	"testing"
	"time"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_EditChatMessage(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	messageID := uuid.New()
	ctx := auth.WithUserID(context.Background(), userID)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		editedAt := time.Now()
		expected := &models.ChatMessage{ID: messageID, SenderID: userID, Message: "Edited", EditedAt: &editedAt}
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("EditChatMessage", ctx, messageID, userID, "Edited").Return(expected, nil)

		res, err := NewHandler(mockRepo).EditChatMessage(ctx, messageID, "Edited")

		require.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("blank message is rejected", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		res, err := NewHandler(mockRepo).EditChatMessage(ctx, messageID, "   ")

		assert.Nil(t, res)
		assertHTTPStatus(t, err, http.StatusBadRequest)
	})

	t.Run("not found is passed through", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("EditChatMessage", ctx, messageID, userID, "Edited").
			Return(nil, errs.NotFound("chat message", "id", messageID))

		res, err := NewHandler(mockRepo).EditChatMessage(ctx, messageID, "Edited")

		assert.Nil(t, res)
		assertHTTPStatus(t, err, http.StatusNotFound)
	})

	t.Run("missing user is unauthorized", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		res, err := NewHandler(mockRepo).EditChatMessage(context.Background(), messageID, "Edited")

		assert.Nil(t, res)
		assertHTTPStatus(t, err, http.StatusUnauthorized)
	})
}

func TestHandler_UpdateChatMessageReadAt(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	messageID := uuid.New()
	ctx := auth.WithUserID(context.Background(), userID)

	t.Run("read=true passes the current time", func(t *testing.T) {
		t.Parallel()

		before := time.Now()
		expected := &models.ChatMessage{ID: messageID, RecipientID: userID}
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("UpdateChatMessageReadAt", ctx, messageID, userID, mock.MatchedBy(func(readAt *time.Time) bool {
			return readAt != nil && !readAt.Before(before) && !readAt.After(time.Now())
		})).Return(expected, nil)

		res, err := NewHandler(mockRepo).UpdateChatMessageReadAt(ctx, messageID, true)

		require.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("read=false passes nil", func(t *testing.T) {
		t.Parallel()

		expected := &models.ChatMessage{ID: messageID, RecipientID: userID}
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("UpdateChatMessageReadAt", ctx, messageID, userID, (*time.Time)(nil)).Return(expected, nil)

		res, err := NewHandler(mockRepo).UpdateChatMessageReadAt(ctx, messageID, false)

		require.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestHandler_MarkChatRead(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	otherID := uuid.New()
	ctx := auth.WithUserID(context.Background(), userID)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("MarkChatRead", ctx, userID, otherID).Return(nil)

		assert.NoError(t, NewHandler(mockRepo).MarkChatRead(ctx, otherID))
	})
}
