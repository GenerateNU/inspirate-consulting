-- Remove students, keeping tasks they already completed.
DELETE FROM public.todo_items
WHERE assignment_id = $1
  AND student_id = ANY($2::uuid[])
  AND completed_at IS NULL