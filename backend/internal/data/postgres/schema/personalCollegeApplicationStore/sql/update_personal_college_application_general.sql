UPDATE public.personal_college_applications
SET
    global_college_id = $3,
    application_type = $4,
    category = $5
WHERE id = $1 AND student_id = $2
RETURNING id, student_id, global_college_id, application_type, category, rank, created_at, updated_at