-- Scoped to the parent so an achievement can't be edited through another
-- extracurricular's URL. NULL parameters mean "not sent" and keep the value.
UPDATE public.extracurricular_achievements SET
    entry_date  = COALESCE($3, entry_date),
    body        = COALESCE($4, body),
    modified_at = now()
WHERE id = $2
  AND extracurricular_id = $1
RETURNING id, extracurricular_id, entry_date, body, created_at, modified_at;
