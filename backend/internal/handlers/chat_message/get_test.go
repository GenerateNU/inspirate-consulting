package chatmessage

import (
	"context"
	"net/http"
	"testing"
	"time"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// newestFirstMessages builds n messages between two users, newest first.
func newestFirstMessages(n int, a, b uuid.UUID) []models.ChatMessage {
	base := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	msgs := make([]models.ChatMessage, n)
	for i := range msgs {
		msgs[i] = models.ChatMessage{
			ID:          uuid.New(),
			SenderID:    a,
			RecipientID: b,
			Message:     "message",
			CreatedAt:   base.Add(-time.Duration(i) * time.Minute),
		}
	}
	return msgs
}

func TestHandler_ListChatMessages(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	otherID := uuid.New()
	ctx := auth.WithUserID(context.Background(), userID)
	nilKey := (*pagination.TimeIDKey)(nil)

	t.Run("first page with no more messages has no cursor", func(t *testing.T) {
		t.Parallel()

		rows := newestFirstMessages(2, otherID, userID)
		mockRepo := mocks.NewChatMessageRepository(t)
		// The handler asks for one extra row to detect another page.
		mockRepo.On("ListChatMessages", ctx, userID, otherID, nilKey, 4).Return(rows, nil)

		page, err := NewHandler(mockRepo).ListChatMessages(ctx, otherID, "", 3)

		require.NoError(t, err)
		assert.Equal(t, rows, page.Messages)
		assert.Nil(t, page.NextCursor)
	})

	t.Run("exactly limit rows has no cursor", func(t *testing.T) {
		t.Parallel()

		rows := newestFirstMessages(3, otherID, userID)
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("ListChatMessages", ctx, userID, otherID, nilKey, 4).Return(rows, nil)

		page, err := NewHandler(mockRepo).ListChatMessages(ctx, otherID, "", 3)

		require.NoError(t, err)
		assert.Len(t, page.Messages, 3)
		assert.Nil(t, page.NextCursor)
	})

	t.Run("extra row is trimmed and cursor points at last returned message", func(t *testing.T) {
		t.Parallel()

		rows := newestFirstMessages(4, otherID, userID)
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("ListChatMessages", ctx, userID, otherID, nilKey, 4).Return(rows, nil)

		page, err := NewHandler(mockRepo).ListChatMessages(ctx, otherID, "", 3)

		require.NoError(t, err)
		assert.Equal(t, rows[:3], page.Messages)
		require.NotNil(t, page.NextCursor)

		key, err := pagination.Decode[pagination.TimeIDKey](*page.NextCursor)
		require.NoError(t, err)
		assert.Equal(t, rows[2].ID, key.ID)
		assert.True(t, rows[2].CreatedAt.Equal(key.CreatedAt))
	})

	t.Run("cursor is decoded and passed to the repository", func(t *testing.T) {
		t.Parallel()

		want := pagination.TimeIDKey{CreatedAt: time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC), ID: uuid.New()}
		cursor, err := pagination.Encode(want)
		require.NoError(t, err)

		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("ListChatMessages", ctx, userID, otherID, mock.MatchedBy(func(k *pagination.TimeIDKey) bool {
			return k != nil && k.ID == want.ID && k.CreatedAt.Equal(want.CreatedAt)
		}), 11).Return([]models.ChatMessage{}, nil)

		page, err := NewHandler(mockRepo).ListChatMessages(ctx, otherID, cursor, 10)

		require.NoError(t, err)
		assert.Empty(t, page.Messages)
		assert.Nil(t, page.NextCursor)
	})

	t.Run("invalid cursor is rejected", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewChatMessageRepository(t)
		page, err := NewHandler(mockRepo).ListChatMessages(ctx, otherID, "not-a-valid-cursor!", 10)

		assert.Nil(t, page)
		assertHTTPStatus(t, err, http.StatusBadRequest)
	})
}
