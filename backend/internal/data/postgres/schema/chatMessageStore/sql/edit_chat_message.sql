UPDATE public.chat_messages
SET message = $1, edited_at = now()
WHERE id = $2 AND sender_id = $3
RETURNING id, sender_id, recipient_id, message, created_at, edited_at, read_at
