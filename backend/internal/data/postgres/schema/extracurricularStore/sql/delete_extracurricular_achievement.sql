DELETE FROM public.extracurricular_achievements
WHERE id = $2
  AND extracurricular_id = $1;
