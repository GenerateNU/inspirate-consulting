-- Idempotent: COALESCE keeps the first timestamp on repeat calls.
-- Completing also marks the media started, in case the started call never happened.
UPDATE public.media_access
SET started_at = COALESCE(started_at, now()),
    completed_at = CASE WHEN $3 THEN COALESCE(completed_at, now()) ELSE completed_at END
WHERE student_id = $1 AND media_id = $2
RETURNING id, student_id, media_id, started_at, completed_at;
