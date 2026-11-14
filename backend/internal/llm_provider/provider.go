package llmprovider

import (
	"context"

	"inspirate-consulting/internal/models"
)

type LLM string

const (
	GeminiModel LLM = "gemini"
)

// Searcher runs the similarity searches behind the search_chats and search_tasks tools
type Searcher interface {
	SearchChats(ctx context.Context, query string) ([]models.ChatMessageSearchResult, error)
	SearchTasks(ctx context.Context, query string) ([]models.TodoItemSearchResult, error)
}

type LLMProvider struct {
	UserInput string
	Model     LLM
	Search    Searcher
}
