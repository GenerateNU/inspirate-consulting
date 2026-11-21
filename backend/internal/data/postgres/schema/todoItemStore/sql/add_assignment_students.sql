-- Add students, copying the assignment's current fields.
-- Only the counselor's own students, and only ones not already in it.
INSERT INTO public.todo_items
    (user_id, title, todo_description, deadline, link_type,
     global_college_id, media_id, assignment_id, student_id)
SELECT t.user_id, t.title, t.todo_description, t.deadline, t.link_type,
       t.global_college_id, t.media_id, t.assignment_id, s.id
FROM (
    SELECT * FROM public.todo_items WHERE assignment_id = $1 LIMIT 1
) t
CROSS JOIN public.student s
WHERE s.id = ANY($3::uuid[])
  AND s.counselor_id = $2
  AND NOT EXISTS (
      SELECT 1 FROM public.todo_items x
      WHERE x.assignment_id = $1 AND x.student_id = s.id
  )