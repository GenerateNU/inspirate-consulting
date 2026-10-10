SELECT id, rank
FROM public.personal_college_applications
WHERE student_id = $1
ORDER BY rank NULLS LAST
FOR UPDATE