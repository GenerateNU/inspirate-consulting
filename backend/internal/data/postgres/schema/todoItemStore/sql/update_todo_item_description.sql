UPDATE public.todo_items
SET todo_description = $2, updated_at = now()
WHERE id = $1
RETURNING id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline, essay_id, media_id, global_college_id