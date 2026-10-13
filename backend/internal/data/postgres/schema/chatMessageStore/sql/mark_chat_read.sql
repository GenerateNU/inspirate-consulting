UPDATE public.chat_messages
SET read_at = now()
WHERE recipient_id = $1 AND sender_id = $2 AND read_at IS NULL
