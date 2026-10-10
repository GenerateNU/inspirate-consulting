WITH filtered AS (
    SELECT id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline, essay_id, media_id, global_college_id,
        -- Missing dates are replaced with $10 so the sort key is never null and can be compared against the cursor
        coalesce(
            CASE $8::text
                WHEN 'deadline' THEN deadline
                WHEN 'completed_at' THEN completed_at
                WHEN 'updated_at' THEN updated_at
                ELSE created_at
            END,
            $10::timestamptz
        ) AS sort_key
    FROM public.todo_items
    WHERE student_id = $1
      AND (
        $2::text = 'all'
        OR ($2 = 'completed' AND completed_at IS NOT NULL)
        OR ($2 = 'incomplete' AND completed_at IS NULL)
      )
      AND ($3::uuid IS NULL OR essay_id = $3)
      AND ($4::uuid IS NULL OR media_id = $4)
      AND ($5::bigint IS NULL OR global_college_id = $5)
      AND (
        $6::text IS NULL
        OR ($6 = 'essay' AND essay_id IS NOT NULL)
        OR ($6 = 'media' AND media_id IS NOT NULL)
        OR ($6 = 'college' AND global_college_id IS NOT NULL)
      )
      AND ($7::text IS NULL OR todo_description ILIKE '%' || $7 || '%')
)
SELECT id, created_at, updated_at, student_id, user_id, todo_description, completed_at, deadline, essay_id, media_id, global_college_id
FROM filtered
WHERE $11::timestamptz IS NULL
   OR ($9::text = 'asc' AND (sort_key, id) > ($11, $12::uuid))
   OR ($9 = 'desc' AND (sort_key, id) < ($11, $12::uuid))
ORDER BY
    CASE WHEN $9 = 'asc' THEN sort_key END ASC,
    CASE WHEN $9 = 'asc' THEN id END ASC,
    CASE WHEN $9 = 'desc' THEN sort_key END DESC,
    CASE WHEN $9 = 'desc' THEN id END DESC
LIMIT $13
