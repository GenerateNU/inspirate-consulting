package llmprovider

import (
	"context"
	"errors"
	"testing"
	"time"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchChats(t *testing.T) {
	t.Parallel()

	me := uuid.New()
	student := uuid.New()
	ctx := auth.WithUserID(context.Background(), me)

	t.Run("lists matches with sender and date", func(t *testing.T) {
		t.Parallel()

		searcher := &fakeSearcher{chats: []models.ChatMessageSearchResult{
			{ChatMessage: models.ChatMessage{SenderID: student, RecipientID: me, Message: "I hate literature class", CreatedAt: time.Date(2026, 11, 3, 15, 0, 0, 0, time.UTC)}, Similarity: 0.9},
			{ChatMessage: models.ChatMessage{SenderID: me, RecipientID: student, Message: "What part is hardest?", CreatedAt: time.Date(2026, 11, 4, 9, 30, 0, 0, time.UTC)}, Similarity: 0.6},
		}}
		provider := LLMProvider{Search: searcher}

		answer, err := provider.searchChats(ctx, "student complaining about literature")

		require.NoError(t, err)
		require.NotNil(t, answer)
		assert.Equal(t,
			"- Them (Nov 3, 2026): I hate literature class\n"+
				"- You (Nov 4, 2026): What part is hardest?\n",
			*answer)
		assert.Equal(t, "student complaining about literature", searcher.gotQuery)
	})

	t.Run("says when nothing matches", func(t *testing.T) {
		t.Parallel()

		provider := LLMProvider{Search: &fakeSearcher{chats: []models.ChatMessageSearchResult{}}}

		answer, err := provider.searchChats(ctx, "anything")

		require.NoError(t, err)
		require.NotNil(t, answer)
		assert.Equal(t, "No matching messages found.", *answer)
	})

	t.Run("search error is returned", func(t *testing.T) {
		t.Parallel()

		searchErr := errors.New("similarity search is not configured")
		provider := LLMProvider{Search: &fakeSearcher{err: searchErr}}

		answer, err := provider.searchChats(ctx, "anything")

		assert.Nil(t, answer)
		assert.ErrorIs(t, err, searchErr)
	})

	t.Run("errors without a searcher", func(t *testing.T) {
		t.Parallel()

		provider := LLMProvider{}

		answer, err := provider.searchChats(ctx, "anything")

		assert.Nil(t, answer)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not configured")
	})
}
