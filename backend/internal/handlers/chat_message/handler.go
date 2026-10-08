package chatmessage

import (
	"context"

	"inspirate-consulting/internal/auth"
	storage "inspirate-consulting/internal/data"
	"inspirate-consulting/internal/errs"

	"github.com/google/uuid"
)

type Handler struct {
	ChatMessageRepository storage.ChatMessageRepository
}

// Handler with reference to ChatMessageRepository
func NewHandler(chatMessageRepository storage.ChatMessageRepository) *Handler {
	return &Handler{
		ChatMessageRepository: chatMessageRepository,
	}
}

func currentUserID(ctx context.Context) (uuid.UUID, error) {
	id, err := uuid.Parse(auth.GetUserID(ctx))
	if err != nil {
		return uuid.Nil, errs.Unauthorized()
	}
	return id, nil
}
