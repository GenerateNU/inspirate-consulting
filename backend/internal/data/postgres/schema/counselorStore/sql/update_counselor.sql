-- STUB: the counselor table has no profile fields yet
UPDATE public.counselor
SET
    updated_at = NOW()
WHERE id = $1
RETURNING id, user_id
