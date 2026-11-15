UPDATE public.personal_college_applications
SET
    global_college_id = $3,
    application_type = $4,
    category = $5,
    interest = $6,
    application_status = $7,
    updated_at = now()
WHERE id = $1 AND student_id = $2
RETURNING id
