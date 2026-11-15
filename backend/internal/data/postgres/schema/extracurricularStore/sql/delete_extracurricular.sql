-- Achievements and notes are removed by ON DELETE CASCADE.
DELETE FROM public.extracurriculars
WHERE id = $1
  AND student_id = $2;
