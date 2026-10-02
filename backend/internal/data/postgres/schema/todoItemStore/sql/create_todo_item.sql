INSERT INTO public.todo_items (
    student_id, user_id, todo_description, deadline
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline