INSERT INTO public.notification_preferences (user_id)
VALUES ($1)
RETURNING user_id, email_enabled, weekly_summary_enabled, due_date_notifications_enabled, days_before_due, notify_past_due