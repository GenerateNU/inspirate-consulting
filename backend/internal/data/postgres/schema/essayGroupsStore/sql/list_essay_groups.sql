SELECT id, name, description, student_id
FROM public.essay_groups
WHERE student_id = $1
