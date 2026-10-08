package chatmessage

import (
	"context"

	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"

	"github.com/google/uuid"
)

func (h *Handler) ListChatMessages(ctx context.Context, otherUserID uuid.UUID, cursor string, limit int) (*models.ChatMessagePage, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	var before *pagination.TimeIDKey
	if cursor != "" {
		key, err := pagination.Decode[pagination.TimeIDKey](cursor)
		if err != nil {
			return nil, err
		}
		before = &key
	}

	rows, err := h.ChatMessageRepository.ListChatMessages(ctx, userID, otherUserID, before, limit+1)
	if err != nil {
		return nil, err
	}

	messages, nextCursor, err := pagination.NextPage(rows, limit, func(m models.ChatMessage) pagination.TimeIDKey {
		return pagination.TimeIDKey{CreatedAt: m.CreatedAt, ID: m.ID}
	})
	if err != nil {
		return nil, err
	}

	return &models.ChatMessagePage{Messages: messages, NextCursor: nextCursor}, nil
}
