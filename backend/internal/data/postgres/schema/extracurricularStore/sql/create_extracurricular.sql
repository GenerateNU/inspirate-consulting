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
    organization,
    hours_per_week,
    weeks_per_year,
    goals,
    ready_for_application,
    application_description
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, COALESCE($14, false), $15
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
    organization,
    hours_per_week,
    weeks_per_year,
    goals,
    ready_for_application,
    application_description;
