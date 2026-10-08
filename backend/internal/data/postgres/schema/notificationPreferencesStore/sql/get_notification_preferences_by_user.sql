SELECT user_id, email_enabled, weekly_summary_enabled, due_date_notifications_enabled, days_before_due, notify_past_due
FROM public.notification_preferences
WHERE user_id = $1