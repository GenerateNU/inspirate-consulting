DELETE FROM public.todo_items
WHERE id = $1
RETURNING student_id, media_id
