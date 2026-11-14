-- Scoped to the parent so a note can't be edited through another extracurricular's URL.
UPDATE public.extracurricular_notes SET
    body        = $3,
    modified_at = now()
WHERE id = $2
  AND extracurricular_id = $1
RETURNING id, extracurricular_id, body, created_at, modified_at;
