package models

// A chat message matched by similarity search, for the search_chats LLM tool.
type ChatMessageSearchResult struct {
	ChatMessage
	Similarity float64 `json:"similarity" doc:"Cosine similarity to the search query, 1 = identical meaning"`
}

// A todo item matched by similarity search, for the search_tasks LLM tool.
type TodoItemSearchResult struct {
	TodoItem
	Similarity float64 `json:"similarity" doc:"Cosine similarity to the search query, 1 = identical meaning"`
}
