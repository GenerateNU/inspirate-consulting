package notificationpreferences

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) UpdateNotificationPreferences(ctx context.Context, input models.UpdateNotificationPreferencesRequestBody) (*models.NotificationPreferences, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	updatedNotificationPreferences, err := h.NotificationPreferencesRepository.UpdateNotificationPreferences(ctx, userID, input)
	if err != nil {
		return nil, err
	}

	return updatedNotificationPreferences, nil
}
