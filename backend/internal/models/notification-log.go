package models

import "time"

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

type TaskNotification struct {
	TaskID          string    `json:"task_id"`
	Deadline        time.Time `json:"deadline"`
	TodoDescription string    `json:"todo_description"`
	UserID          string    `json:"user_id"`
}
