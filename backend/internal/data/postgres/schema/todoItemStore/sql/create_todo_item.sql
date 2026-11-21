INSERT INTO public.todo_items (
    student_id, user_id, todo_description, deadline, essay_id, media_id, global_college_id, title, extracurricular_id, link_type
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::public.todo_link_type
)
RETURNING id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline, essay_id, media_id, global_college_id, title, extracurricular_id, link_type, assignment_id;
