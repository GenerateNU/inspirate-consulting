package models

import "time"

// enum type matching NOTIFICATION_TYPE in the notification_logs migration
type NotificationType string

const (
	NotificationUpcoming NotificationType = "upcoming"
	NotificationPastDue  NotificationType = "past_due"
)

type NotificationLog struct {
	ID               string           `json:"id"`
	UserID           string           `json:"user_id"`
	TaskID           string           `json:"task_id"`
	NotificationType NotificationType `json:"notification_type"`
	SentAt           time.Time        `json:"sent_at"`
	CreatedAt        time.Time        `json:"created_at"`
}

type CreateNotificationLogInput struct {
	UserID           string           `json:"user_id"`
	TaskID           string           `json:"task_id"`
	NotificationType NotificationType `json:"notification_type"`
	SentAt           time.Time        `json:"sent_at"`
}
