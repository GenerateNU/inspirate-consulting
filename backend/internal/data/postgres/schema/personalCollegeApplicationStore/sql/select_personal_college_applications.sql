-- $1 = student_id, $2 = application id (NULL to list all of the student's applications)
SELECT
    pca.id,
    pca.student_id,
    pca.global_college_id,
    pca.application_type,
    pca.category,
    pca.created_at,
    pca.updated_at,
    pca.rank,
    pca.interest,
    pca.application_status,
    EXISTS (
        SELECT 1
        FROM public.college_application_cycles c
        WHERE c.global_college_id = pca.global_college_id AND c.cycle = pca.application_type
    ) AS cycle_offered,
    CASE
        WHEN ep.total = 0 THEN 'not_started'
        WHEN ep.submitted = ep.total THEN 'completed'
        ELSE 'in_progress'
    END AS essay_progress
FROM public.personal_college_applications pca
CROSS JOIN LATERAL (
    SELECT
        count(*) AS total,
        count(*) FILTER (WHERE e.status = 'Submitted') AS submitted
    FROM public.essays e
    WHERE e.student_id = pca.student_id
      AND e.college_id = pca.global_college_id
      AND e.status <> 'Archived'
) ep
WHERE pca.student_id = $1
  AND ($2::bigint IS NULL OR pca.id = $2)
ORDER BY pca.rank NULLS LAST, pca.id
