package chatMessageRepository

import (
	"context"
	"time"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ListChatMessages returns a page of the conversation between two users, newest first.
func (r *ChatMessageRepository) ListChatMessages(ctx context.Context, userID uuid.UUID, otherUserID uuid.UUID, before *pagination.TimeIDKey, limit int) ([]models.ChatMessage, error) {
	selectQuery, err := schema.ReadSQLBaseScript("list_chat_messages.sql", SqlChatMessageFiles)
	if err != nil {
		return nil, err
	}

	var beforeAt *time.Time
	var beforeID *uuid.UUID
	if before != nil {
		beforeAt, beforeID = &before.CreatedAt, &before.ID
	}

	rows, err := r.db.Query(ctx, selectQuery, userID, otherUserID, beforeAt, beforeID, limit)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.ChatMessage])
}
