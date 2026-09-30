SELECT id, student_id, global_college_id, application_type, category, created_at, updated_at
FROM public.personal_college_applications
WHERE student_id = $1