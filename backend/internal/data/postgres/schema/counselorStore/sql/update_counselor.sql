-- Each nullable column takes a (sent boolean, value) pair; when sent is true the
-- value is written as-is, so NULL clears it. expires_at lives on users.
WITH c AS (
    UPDATE public.counselor SET
        custom_link = CASE WHEN $2::boolean THEN $3::text ELSE custom_link END,
        updated_at  = NOW()
    WHERE id = $1
    RETURNING id, user_id, custom_link
), u AS (
    UPDATE public.users SET
        expires_at = CASE WHEN $4::boolean THEN $5::timestamptz ELSE expires_at END,
        updated_at = NOW()
    WHERE id = (SELECT user_id FROM c)
    RETURNING expires_at
)
SELECT c.id, c.user_id, c.custom_link, u.expires_at
FROM c, u
