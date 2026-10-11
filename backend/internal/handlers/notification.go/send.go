package notification

import (
	"context"
	"inspirate-consulting/internal/models"
)

func (h *Handler) QueueEmailNotification(ctx context.Context, body *models.QueueEmailNotificationRequestBody)(*string, error){
	messageId, err := h.SQSClient.SendEmailNotification(ctx, body)
	if(err != nil){
		return nil, err
	}

	return messageId, nil
}
