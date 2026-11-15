package vectorsearch

import storage "inspirate-consulting/internal/data"

// How many results each search returns.
const searchLimit = 10

// Similarity search over chat messages and todo items, for the LLM's search_chats and search_tasks tools.
type Handler struct {
	VectorRepository    storage.VectorRepository
	CounselorRepository storage.CounselorRepository
	// nil when Cloudflare isn't configured; searches then return an error
	EmbeddingClient storage.EmbeddingClient
}

func NewHandler(vectorRepository storage.VectorRepository, counselorRepository storage.CounselorRepository, embeddingClient storage.EmbeddingClient) *Handler {
	return &Handler{
		VectorRepository:    vectorRepository,
		CounselorRepository: counselorRepository,
		EmbeddingClient:     embeddingClient,
	}
}
