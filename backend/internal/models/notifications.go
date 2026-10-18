package models

import "time"

type QueueEmailNotificationRequestBody struct {
	UserEmail string    `json:"user_email"`
	TaskName  string    `json:"task_name"`
	DueDate   time.Time `json:"due_date"`
	IsOverdue bool      `json:"is_overdue"`
}

type QueueEmailNotificationResponseBody struct {
	MessageID *string `json:"message_id"`
}

type QueueEmailNotificationInput struct {
	Body QueueEmailNotificationRequestBody
}

type QueueEmailNotificationOutput struct {
	Body QueueEmailNotificationResponseBody
}
