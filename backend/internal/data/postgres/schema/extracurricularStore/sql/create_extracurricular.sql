INSERT INTO public.extracurriculars (
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
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING
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
    organization;