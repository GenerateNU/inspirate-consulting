INSERT INTO public.essay_groups (
    name, description, student_id
) VALUES (
    $1, $2, $3
)
RETURNING id, name, description, student_id
