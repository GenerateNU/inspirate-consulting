INSERT INTO public.extracurricular_notes (extracurricular_id, body)
VALUES ($1, $2)
RETURNING id, extracurricular_id, body, created_at, modified_at;
