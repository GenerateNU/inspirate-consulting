INSERT INTO public.todo_items (
<<<<<<< HEAD
    student_id, user_id, todo_description, deadline, essay_id, media_id, global_college_id, title, extracurricular_id, link_type
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10::public.todo_link_type
=======
    student_id, user_id, todo_description, deadline, essay_id, media_id, global_college_id, title
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, left($3, 80)
>>>>>>> main
)
RETURNING id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline, essay_id, media_id, global_college_id, title, extracurricular_id, link_type, assignment_id;
