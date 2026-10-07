package chatMessageRepository

import (
	"context"
	"errors"
	"time"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *ChatMessageRepository) EditChatMessage(ctx context.Context, id uuid.UUID, senderID uuid.UUID, message string) (*models.ChatMessage, error) {
	updateQuery, err := schema.ReadSQLBaseScript("edit_chat_message.sql", SqlChatMessageFiles)
	if err != nil {
		return nil, err
	}

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
	updateQuery, err := schema.ReadSQLBaseScript("update_chat_message_read_at.sql", SqlChatMessageFiles)
	if err != nil {
		return nil, err
	}

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
	updateQuery, err := schema.ReadSQLBaseScript("mark_chat_read.sql", SqlChatMessageFiles)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, updateQuery, userID, otherUserID)
	return err
}
