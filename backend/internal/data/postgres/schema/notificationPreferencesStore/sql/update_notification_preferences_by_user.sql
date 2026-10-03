	UPDATE public.notification_preferences
	SET
		email_enabled = $2,
		weekly_summary_enabled = $3,
		due_date_notifications_enabled = $4,
		days_before_due = $5,
		notify_past_due = $6
	WHERE user_id = $1
	RETURNING user_id, email_enabled, weekly_summary_enabled, due_date_notifications_enabled, days_before_due, notify_past_due