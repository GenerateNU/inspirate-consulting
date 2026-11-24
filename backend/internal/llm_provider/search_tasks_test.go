package llmprovider

import (
	"context"
	"errors"
	"testing"
	"time"

	"inspirate-consulting/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchTasks(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("lists matches with due date and completion status", func(t *testing.T) {
		t.Parallel()

		deadline := time.Date(2026, 12, 1, 23, 59, 0, 0, time.UTC)
		completedAt := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
		searcher := &fakeSearcher{tasks: []models.TodoItemSearchResult{
			{TodoItem: models.TodoItem{TodoDescription: "Finish Common App essay", Deadline: &deadline}, Similarity: 0.9},
			{TodoItem: models.TodoItem{TodoDescription: "Draft activities list", Deadline: &deadline, CompletedAt: &completedAt}, Similarity: 0.7},
			{TodoItem: models.TodoItem{TodoDescription: "Ask for recommendation letters"}, Similarity: 0.5},
		}}
		provider := LLMProvider{Search: searcher}

		answer, err := provider.searchTasks(ctx, "essay deadlines")

		require.NoError(t, err)
		require.NotNil(t, answer)
		assert.Equal(t,
			"- Finish Common App essay (due Dec 1, 2026, not completed)\n"+
				"- Draft activities list (due Dec 1, 2026, completed)\n"+
				"- Ask for recommendation letters (no due date, not completed)\n",
			*answer)
		assert.Equal(t, "essay deadlines", searcher.gotQuery)
	})

	t.Run("says when nothing matches", func(t *testing.T) {
		t.Parallel()

		provider := LLMProvider{Search: &fakeSearcher{tasks: []models.TodoItemSearchResult{}}}

		answer, err := provider.searchTasks(ctx, "anything")

		require.NoError(t, err)
		require.NotNil(t, answer)
		assert.Equal(t, "No matching tasks found.", *answer)
	})

	t.Run("search error is returned", func(t *testing.T) {
		t.Parallel()

		searchErr := errors.New("insufficient role")
		provider := LLMProvider{Search: &fakeSearcher{err: searchErr}}

		answer, err := provider.searchTasks(ctx, "anything")

		assert.Nil(t, answer)
		assert.ErrorIs(t, err, searchErr)
	})

	t.Run("errors without a searcher", func(t *testing.T) {
		t.Parallel()

		provider := LLMProvider{}

		answer, err := provider.searchTasks(ctx, "anything")

		assert.Nil(t, answer)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not configured")
	})
}
