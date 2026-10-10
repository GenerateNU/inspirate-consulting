SELECT m.id, m.title, m.description, m.length_in_mins, m.school_year, m.s3_key
FROM public.media m
JOIN public.media_access ma ON ma.media_id = m.id
WHERE ma.student_id = $1
ORDER BY m.title, m.id
LIMIT $2 OFFSET $3;