package llmprovider

import (
	"context"

	"inspirate-consulting/internal/models"
)

// fakeSearcher returns fixed results (or an error) and records the query it was given.
type fakeSearcher struct {
	chats    []models.ChatMessageSearchResult
	tasks    []models.TodoItemSearchResult
	err      error
	gotQuery string
}

func (f *fakeSearcher) SearchChats(ctx context.Context, query string) ([]models.ChatMessageSearchResult, error) {
	f.gotQuery = query
	return f.chats, f.err
}

func (f *fakeSearcher) SearchTasks(ctx context.Context, query string) ([]models.TodoItemSearchResult, error) {
	f.gotQuery = query
	return f.tasks, f.err
}
