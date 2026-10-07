package chatMessageRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatMessageRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlChatMessageFiles embed.FS

func NewChatMessageRepository(db *pgxpool.Pool) *ChatMessageRepository {
	return &ChatMessageRepository{db: db}
}
