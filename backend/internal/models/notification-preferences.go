package models

import (
	"github.com/google/uuid"
)

// To represent a user's notification preferences.
type NotificationPreferences struct {
	UserID uuid.UUID `json:"user_id" doc:"ID of the user who owns this notification preference"`
	EmailEnabled bool `json:"email_enabled" doc:"Whether the user wants to receive email notifications"`
	WeeklySummaryEnabled bool `json:"weekly_summary_enabled" doc:"Whether the user wants to receive weekly summary notifications"`
	DueDateNotificationsEnabled bool `json:"due_date_notifications_enabled" doc:"Whether the user wants to receive due date notifications"`
	DaysBeforeDue int `json:"days_before_due" doc:"Number of days before the due date to send notifications"`
	NotifyPastDue bool `json:"notify_past_due" doc:"Whether the user wants to receive notifications for past due dates"`
}

type UpdateNotificationPreferencesRequestBody struct {
	EmailEnabled bool `json:"email_enabled" doc:"Whether the user wants to receive email notifications"`
	WeeklySummaryEnabled bool `json:"weekly_summary_enabled" doc:"Whether the user wants to receive weekly summary notifications"`
	DueDateNotificationsEnabled bool `json:"due_date_notifications_enabled" doc:"Whether the user wants to receive due date notifications"`
	DaysBeforeDue int `json:"days_before_due" doc:"Number of days before the due date to send notifications"`
	NotifyPastDue bool `json:"notify_past_due" doc:"Whether the user wants to receive notifications for past due dates"`
}

// Huma readable input and output models, user ID is derived from auth

type GetNotificationPreferencesInput struct{}
 
type GetNotificationPreferencesOutput struct {
	Body NotificationPreferences
}
 
type UpdateNotificationPreferencesInput struct {
	Body UpdateNotificationPreferencesRequestBody
}
 
type UpdateNotificationPreferencesOutput struct {
	Body NotificationPreferences
}

