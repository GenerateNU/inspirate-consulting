DELETE FROM public.extracurricular_notes
WHERE id = $2
  AND extracurricular_id = $1;
