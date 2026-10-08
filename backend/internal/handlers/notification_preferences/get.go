package notificationpreferences

import (
	"context"
	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

// to get notification preferences for a given user
// if the user does not have any notification preferences, a default row will be created for them
func (h *Handler) GetNotificationPreferences(ctx context.Context) (*models.NotificationPreferences, error) {

	userID, err := uuid.Parse(auth.GetUserID(ctx))
	if err != nil {
		return nil, errs.InternalServerError("invalid user ID from auth context")
	}

	notificationPreferences, err := h.NotificationPreferencesRepository.GetNotificationPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}

	return notificationPreferences, nil
}
