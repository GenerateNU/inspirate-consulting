INSERT INTO public.chat_messages (sender_id, recipient_id, message)
VALUES ($1, $2, $3)
RETURNING id, sender_id, recipient_id, message, created_at, edited_at, read_at
