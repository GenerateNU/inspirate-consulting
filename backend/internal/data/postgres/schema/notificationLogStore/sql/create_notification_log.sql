INSERT INTO public.notification_logs (
    user_id, task_id, notification_type, sent_at, created_at
) VALUES (
    $1, $2, $3, $4, now()
)
RETURNING id, user_id, task_id, notification_type, sent_at, created_at
