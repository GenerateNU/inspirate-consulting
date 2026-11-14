package activity

import (
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	ActivityRepository storage.ActivityRepository
}

// Handler with reference to ActivityRepository
func NewHandler(activityRepository storage.ActivityRepository) *Handler {
	return &Handler{
		ActivityRepository: activityRepository,
	}
}
