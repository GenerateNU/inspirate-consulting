-- One row per student, all sharing one assignment id.
INSERT INTO public.todo_items
    (user_id, title, todo_description, deadline, link_type,
     global_college_id, media_id, assignment_id, student_id)
SELECT $1, $2, $3, $4, $5::public.todo_link_type, $6, $7::uuid, $8, sid
FROM unnest($9::uuid[]) AS sid
RETURNING id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline, essay_id, media_id, global_college_id, title, extracurricular_id, link_type, assignment_id