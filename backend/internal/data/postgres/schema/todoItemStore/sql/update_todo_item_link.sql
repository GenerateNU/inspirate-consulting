-- Sets the link type and every link id at once. Whichever ids aren't sent become NULL,
-- so the old link is cleared in the same update and the link checks always pass.
UPDATE public.todo_items
SET link_type          = $2::public.todo_link_type,
    essay_id           = $3::uuid,
    media_id           = $4::uuid,
    global_college_id  = $5,
    extracurricular_id = $6,
    updated_at         = now()
WHERE id = $1
RETURNING id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline, essay_id, media_id, global_college_id, title, extracurricular_id, link_type, assignment_id