package notificationpreferences

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) GetNotificationPreferences(ctx context.Context) (*models.NotificationPreferences, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	notificationPreferences, err := h.NotificationPreferencesRepository.GetNotificationPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}

	return notificationPreferences, nil
}
