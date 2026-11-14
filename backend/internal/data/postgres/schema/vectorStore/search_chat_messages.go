package vectorRepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

// returns the user's messages (sent or received) closest in meaning to the query embedding, most similar first.
func (r *VectorRepository) SearchChatMessages(ctx context.Context, queryEmbedding []float32, userID uuid.UUID, limit int) ([]models.ChatMessageSearchResult, error) {
	searchQuery, err := schema.ReadSQLBaseScript("search_chat_messages.sql", SqlVectorFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, searchQuery, formatVector(queryEmbedding), userID, limit)
	if err != nil {
		return nil, err
	}

	similarChatMessages, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.ChatMessageSearchResult])
	if err != nil {
		return nil, err
	}

	return similarChatMessages, nil

}
