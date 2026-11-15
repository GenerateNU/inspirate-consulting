INSERT INTO public.activity (
    student_id, action, essay_id, essay_status, global_college_id, extracurricular_id, todo_item_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
