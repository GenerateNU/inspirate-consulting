SELECT id, student_id, type, college_id, link_to_content, status, essay_group_id
FROM public.essays
WHERE student_id = $1 AND ($2::bigint IS NULL OR college_id = $2)