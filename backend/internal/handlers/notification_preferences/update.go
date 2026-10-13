package notificationpreferences

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

// to update an existing database entry for a given user
func (h *Handler) UpdateNotificationPreferences(ctx context.Context, input models.UpdateNotificationPreferencesRequestBody) (*models.NotificationPreferences, error) {

	userID, err := uuid.Parse(auth.GetUserID(ctx))
	if err != nil {
		return nil, errs.InternalServerError("invalid user ID from auth context")
	}

	updatedNotificationPreferences, err := h.NotificationPreferencesRepository.UpdateNotificationPreferences(ctx, userID, input)
	if err != nil {
		return nil, err
	}

	return updatedNotificationPreferences, nil
}
