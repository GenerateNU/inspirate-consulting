package chatMessageRepository

import (
	"context"
	"errors"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	pgForeignKeyViolationCode = "23503"
	recipientForeignKey       = "chat_messages_recipient_id_fkey"
)

func (r *ChatMessageRepository) CreateChatMessage(ctx context.Context, senderID uuid.UUID, body models.CreateChatMessageRequestBody) (*models.ChatMessage, error) {
	const insertQuery = `
	INSERT INTO public.chat_messages (sender_id, recipient_id, message)
	VALUES ($1, $2, $3)
	RETURNING ` + chatMessageColumns

	rows, err := r.db.Query(ctx, insertQuery, senderID, body.RecipientID, body.Message)
	if err != nil {
		return nil, err
	}

	created, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[models.ChatMessage])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgForeignKeyViolationCode && pgErr.ConstraintName == recipientForeignKey {
			return nil, errs.NotFound("user", "id", body.RecipientID)
		}
		return nil, err
	}

	return &created, nil
}
