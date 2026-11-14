-- $1 = query embedding, $2 = student ids the caller can see, $3 = limit
SELECT
    ti.id,
    ti.created_at,
    ti.updated_at,
    ti.student_id,
    ti.user_id,
    ti.todo_description,
    ti.completed_at,
    ti.deadline,
    ti.essay_id,
    ti.media_id,
    ti.global_college_id,
    1 - (vt.embedding <=> $1::extensions.vector) AS similarity
FROM public.vector_table vt
JOIN public.vector_relationship vr ON vr.vector_id = vt.vector_id
JOIN public.todo_items ti ON ti.id = vr.todo_id
WHERE ti.student_id = ANY($2::uuid[])
ORDER BY vt.embedding <=> $1::extensions.vector
LIMIT $3
