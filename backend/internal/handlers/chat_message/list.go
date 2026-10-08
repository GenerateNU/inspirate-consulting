package chatmessage

import (
	"context"

	"inspirate-consulting/internal/models"
)

func (h *Handler) ListChats(ctx context.Context) ([]models.ChatSummary, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	return h.ChatMessageRepository.ListChats(ctx, userID)
}
