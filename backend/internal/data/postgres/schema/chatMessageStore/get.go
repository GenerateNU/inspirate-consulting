package chatMessageRepository

import (
	"context"
	"time"

	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ListChatMessages returns a page of the conversation between two users, newest first.
func (r *ChatMessageRepository) ListChatMessages(ctx context.Context, userID uuid.UUID, otherUserID uuid.UUID, before *pagination.TimeIDKey, limit int) ([]models.ChatMessage, error) {
	const selectQuery = `
	SELECT ` + chatMessageColumns + `
	FROM public.chat_messages
	WHERE user_low = least($1::uuid, $2::uuid)
		AND user_high = greatest($1::uuid, $2::uuid)
		AND ($3::timestamptz IS NULL OR (created_at, id) < ($3::timestamptz, $4::uuid))
	ORDER BY created_at DESC, id DESC
	LIMIT $5
	`

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
