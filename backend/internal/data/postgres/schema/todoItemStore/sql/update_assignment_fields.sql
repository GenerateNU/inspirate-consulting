-- Edit every row. When link_type is sent, the link ids are replaced together,
-- so the link checks always pass.
UPDATE public.todo_items
SET title             = COALESCE($2, title),
    todo_description  = COALESCE($3, todo_description),
    deadline          = COALESCE($4, deadline),
    link_type         = COALESCE($5::public.todo_link_type, link_type),
    global_college_id = CASE WHEN $5::text IS NULL THEN global_college_id ELSE $6::bigint END,
    media_id          = CASE WHEN $5::text IS NULL THEN media_id ELSE $7::uuid END
WHERE assignment_id = $1