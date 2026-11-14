INSERT INTO public.extracurricular_achievements (extracurricular_id, entry_date, body)
VALUES ($1, $2, $3)
RETURNING id, extracurricular_id, entry_date, body, created_at, modified_at;
