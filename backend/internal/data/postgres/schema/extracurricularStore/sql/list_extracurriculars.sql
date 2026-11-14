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
    organization,
    hours_per_week,
    weeks_per_year,
    goals,
    ready_for_application,
    application_description
FROM public.extracurriculars
WHERE student_id = $1
ORDER BY start_date DESC;
