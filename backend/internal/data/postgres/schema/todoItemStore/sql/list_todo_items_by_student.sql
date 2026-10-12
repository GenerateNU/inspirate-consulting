SELECT id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline, essay_id, media_id, global_college_id
FROM public.todo_items
WHERE student_id = $1
