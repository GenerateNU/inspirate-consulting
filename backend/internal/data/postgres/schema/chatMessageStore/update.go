package chatMessageRepository

import (
	"context"
	"errors"
	"time"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *ChatMessageRepository) EditChatMessage(ctx context.Context, id uuid.UUID, senderID uuid.UUID, message string) (*models.ChatMessage, error) {
	const updateQuery = `
	UPDATE public.chat_messages
	SET message = $1, edited_at = now()
	WHERE id = $2 AND sender_id = $3
	RETURNING ` + chatMessageColumns

	rows, err := r.db.Query(ctx, updateQuery, message, id, senderID)
	if err != nil {
		return nil, err
	}

	edited, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[models.ChatMessage])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("chat message", "id", id)
		}
		return nil, err
	}

	return &edited, nil
}

func (r *ChatMessageRepository) UpdateChatMessageReadAt(ctx context.Context, id uuid.UUID, recipientID uuid.UUID, readAt *time.Time) (*models.ChatMessage, error) {
	const updateQuery = `
	UPDATE public.chat_messages
	SET read_at = CASE WHEN $1::timestamptz IS NULL THEN NULL ELSE coalesce(read_at, $1::timestamptz) END
	WHERE id = $2 AND recipient_id = $3
	RETURNING ` + chatMessageColumns

	rows, err := r.db.Query(ctx, updateQuery, readAt, id, recipientID)
	if err != nil {
		return nil, err
	}

	updated, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[models.ChatMessage])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("chat message", "id", id)
		}
		return nil, err
	}

	return &updated, nil
}

// MarkChatRead marks every unread message the other user sent to this user as read.
func (r *ChatMessageRepository) MarkChatRead(ctx context.Context, userID uuid.UUID, otherUserID uuid.UUID) error {
	const updateQuery = `
	UPDATE public.chat_messages
	SET read_at = now()
	WHERE recipient_id = $1 AND sender_id = $2 AND read_at IS NULL
	`

	_, err := r.db.Exec(ctx, updateQuery, userID, otherUserID)
	return err
}
