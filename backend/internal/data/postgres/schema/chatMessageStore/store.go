package chatMessageRepository

import "github.com/jackc/pgx/v5/pgxpool"

type ChatMessageRepository struct {
	db *pgxpool.Pool
}

func NewChatMessageRepository(db *pgxpool.Pool) *ChatMessageRepository {
	return &ChatMessageRepository{db: db}
}

const chatMessageColumns = `id, sender_id, recipient_id, message, created_at, edited_at, read_at`
