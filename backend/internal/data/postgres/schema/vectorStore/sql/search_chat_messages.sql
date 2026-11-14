-- $1 = query embedding, $2 = user id (sender or recipient), $3 = limit
SELECT
    cm.id,
    cm.sender_id,
    cm.recipient_id,
    cm.message,
    cm.created_at,
    cm.edited_at,
    cm.read_at,
    1 - (vt.embedding <=> $1::extensions.vector) AS similarity
FROM public.vector_table vt
JOIN public.vector_relationship vr ON vr.vector_id = vt.vector_id
JOIN public.chat_messages cm ON cm.id = vr.message_id
WHERE cm.sender_id = $2 OR cm.recipient_id = $2
ORDER BY vt.embedding <=> $1::extensions.vector
LIMIT $3
