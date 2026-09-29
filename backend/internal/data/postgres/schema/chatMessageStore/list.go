package chatMessageRepository

import (
	"context"

	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ListChats returns one row per conversation the user is in, most recent conversation first.
func (r *ChatMessageRepository) ListChats(ctx context.Context, userID uuid.UUID) ([]models.ChatSummary, error) {
	const selectQuery = `
	WITH latest AS (
		SELECT DISTINCT ON (user_low, user_high)
			id, sender_id, message, created_at, read_at,
			CASE WHEN user_low = $1 THEN user_high ELSE user_low END AS other_user_id
		FROM public.chat_messages
		WHERE user_low = $1 OR user_high = $1
		ORDER BY user_low, user_high, created_at DESC, id DESC
	),
	unread AS (
		SELECT sender_id AS other_user_id, count(*) AS unread_count
		FROM public.chat_messages
		WHERE recipient_id = $1 AND read_at IS NULL
		GROUP BY sender_id
	)
	SELECT
		l.other_user_id, u.name, u.pfp_key,
		l.id, l.message, l.sender_id, l.created_at, l.read_at,
		coalesce(un.unread_count, 0)
	FROM latest l
	JOIN public.users u ON u.id = l.other_user_id
	LEFT JOIN unread un ON un.other_user_id = l.other_user_id
	ORDER BY l.created_at DESC, l.id DESC
	`

	rows, err := r.db.Query(ctx, selectQuery, userID)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[models.ChatSummary])
}
