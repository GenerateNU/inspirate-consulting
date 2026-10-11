INSERT INTO public.personal_college_applications (
    student_id, global_college_id, application_type, category
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, student_id, global_college_id, application_type, category, created_at, updated_at, rank