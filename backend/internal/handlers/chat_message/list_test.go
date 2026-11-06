package chatmessage

import (
	"context"
	"testing"

	"inspirate-consulting/internal/auth"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_ListChats(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	ctx := auth.WithUserID(context.Background(), userID)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		expected := []models.ChatSummary{
			{OtherUserID: uuid.New(), OtherUserName: "Alice", LastMessage: "See you", UnreadCount: 2},
			{OtherUserID: uuid.New(), OtherUserName: "Bob", LastMessage: "Thanks", UnreadCount: 0},
		}
		mockRepo := mocks.NewChatMessageRepository(t)
		mockRepo.On("ListChats", ctx, userID).Return(expected, nil)

		res, err := NewHandler(mockRepo).ListChats(ctx)

		require.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}
