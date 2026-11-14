package vectorsearch

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// finds the current user's chat messages (sent or received) closest in meaning to the query.
func (h *Handler) SearchChats(ctx context.Context, query string) ([]models.ChatMessageSearchResult, error) {
	if h.EmbeddingClient == nil {
		return nil, errs.NewHTTPError(http.StatusServiceUnavailable, errors.New("similarity search is not configured"))
	}
	if strings.TrimSpace(query) == "" {
		return nil, errs.BadRequest("query cannot be blank")
	}
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	queryEmbedding, err := h.EmbeddingClient.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	return h.VectorRepository.SearchChatMessages(ctx, queryEmbedding, userID, searchLimit)
}
