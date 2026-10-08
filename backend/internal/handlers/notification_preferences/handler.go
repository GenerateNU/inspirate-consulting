package notificationpreferences

import storage "inspirate-consulting/internal/data"

type Handler struct {
	NotificationPreferencesRepository storage.NotificationPreferencesRepository
}

// Here is where we create the handler which has reference to the NotificationPreferencesRepository
func NewHandler(notificationPreferencesRepository storage.NotificationPreferencesRepository) *Handler {
	return &Handler{
		NotificationPreferencesRepository: notificationPreferencesRepository,
	}
}