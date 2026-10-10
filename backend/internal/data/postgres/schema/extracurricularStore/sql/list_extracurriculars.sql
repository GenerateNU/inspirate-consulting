SELECT
    id,
    created_at,
    modified_at,
    student_id,
    user_id,
    name,
    status,
    type,
    description,
    leadership_role,
    start_date,
    end_date,
    organization
FROM public.extracurriculars
WHERE student_id = $1
ORDER BY start_date DESC;
