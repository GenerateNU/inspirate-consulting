UPDATE public.essays
SET essay_group_id = $2
WHERE id = $1
RETURNING id, student_id, type, college_id, link_to_content, status, essay_group_id
