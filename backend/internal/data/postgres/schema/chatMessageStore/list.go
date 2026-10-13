package chatMessageRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ListChats returns one row per conversation the user is in, most recent conversation first.
func (r *ChatMessageRepository) ListChats(ctx context.Context, userID uuid.UUID) ([]models.ChatSummary, error) {
	selectQuery, err := schema.ReadSQLBaseScript("list_chats.sql", SqlChatMessageFiles)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, selectQuery, userID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.ChatSummary])
}
