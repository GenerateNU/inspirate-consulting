package chatmessage

import (
	"context"
	"strings"
	"time"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

func (h *Handler) EditChatMessage(ctx context.Context, id uuid.UUID, message string) (*models.ChatMessage, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(message) == "" {
		return nil, errs.BadRequest("message cannot be blank")
	}

	return h.ChatMessageRepository.EditChatMessage(ctx, id, userID, message)
}

func (h *Handler) UpdateChatMessageReadAt(ctx context.Context, id uuid.UUID, read bool) (*models.ChatMessage, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	var readAt *time.Time
	if read {
		now := time.Now()
		readAt = &now
	}

	return h.ChatMessageRepository.UpdateChatMessageReadAt(ctx, id, userID, readAt)
}

func (h *Handler) MarkChatRead(ctx context.Context, otherUserID uuid.UUID) error {
	userID, err := currentUserID(ctx)
	if err != nil {
		return err
	}

	return h.ChatMessageRepository.MarkChatRead(ctx, userID, otherUserID)
}
