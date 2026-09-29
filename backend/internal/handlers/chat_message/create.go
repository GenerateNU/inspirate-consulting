package chatmessage

import (
	"context"
	"strings"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) CreateChatMessage(ctx context.Context, body models.CreateChatMessageRequestBody) (*models.ChatMessage, error) {
	senderID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if body.RecipientID == senderID {
		return nil, errs.BadRequest("cannot send a message to yourself")
	}
	if strings.TrimSpace(body.Message) == "" {
		return nil, errs.BadRequest("message cannot be blank")
	}

	return h.ChatMessageRepository.CreateChatMessage(ctx, senderID, body)
}
