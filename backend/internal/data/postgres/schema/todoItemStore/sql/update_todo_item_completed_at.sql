UPDATE public.todo_items
SET completed_at = $1, updated_at = now()
WHERE id = $2
RETURNING id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline, essay_id, media_id, global_college_id
