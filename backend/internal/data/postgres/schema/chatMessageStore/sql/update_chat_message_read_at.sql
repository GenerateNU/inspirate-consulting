UPDATE public.chat_messages
SET read_at = CASE WHEN $1::timestamptz IS NULL THEN NULL ELSE coalesce(read_at, $1::timestamptz) END
WHERE id = $2 AND recipient_id = $3
RETURNING id, sender_id, recipient_id, message, created_at, edited_at, read_at
