package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	notificationpreferences "inspirate-consulting/internal/handlers/notification_preferences"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

// SetUpNotificationPreferencesRoutes registers the notification preferences endpoints:
// get and update by user ID
func SetUpNotificationPreferencesRoutes(api huma.API, repository *data.Repository) {
	notificationPreferencesHandler := notificationpreferences.NewHandler(repository.NotificationPreferences)

	// Register GET /notification-preferences handler.
	huma.Register(api, huma.Operation{
		OperationID: "get-notification-preferences",
		Method:      http.MethodGet,
		Path:        "/notification-preferences",
		Description: "Get notification preferences for the authenticated user.",
		Tags:        []string{"User Notification Preferences"},
	}, func(ctx context.Context, input *models.GetNotificationPreferencesInput) (*models.GetNotificationPreferencesOutput, error) {
		notificationPreferences, err := notificationPreferencesHandler.GetNotificationPreferences(ctx)
		if err != nil {
			return nil, err
		}
		return &models.GetNotificationPreferencesOutput{Body: *notificationPreferences}, nil
	})

	// Register PUT /notification-preferences handler.
	huma.Register(api, huma.Operation{
		OperationID: "update-notification-preferences",
		Method:      http.MethodPut,
		Path:        "/notification-preferences",
		Description: "Update notification preferences for the authenticated user.",
		Tags:        []string{"User Notification Preferences"},
	}, func(ctx context.Context, input *models.UpdateNotificationPreferencesInput) (*models.UpdateNotificationPreferencesOutput, error) {
		notificationPreferences, err := notificationPreferencesHandler.UpdateNotificationPreferences(ctx, input.Body)
		if err != nil {
			return nil, err
		}
		return &models.UpdateNotificationPreferencesOutput{Body: *notificationPreferences}, nil
	})

}
