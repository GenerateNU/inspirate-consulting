SELECT id, sender_id, recipient_id, message, created_at, edited_at, read_at
FROM public.chat_messages
WHERE user_low = least($1::uuid, $2::uuid)
    AND user_high = greatest($1::uuid, $2::uuid)
    AND ($3::timestamptz IS NULL OR (created_at, id) < ($3::timestamptz, $4::uuid))
ORDER BY created_at DESC, id DESC
LIMIT $5
