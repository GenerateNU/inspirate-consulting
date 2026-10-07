INSERT INTO public.todo_items (
    student_id, user_id, todo_description, deadline, essay_id, media_id, global_college_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline, essay_id, media_id, global_college_id
