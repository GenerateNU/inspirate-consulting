-- "This week" is the current Monday-to-Sunday week in UTC
WITH week AS (
    SELECT date_trunc('week', now() AT TIME ZONE 'UTC') AT TIME ZONE 'UTC' AS week_start
)
SELECT
    COUNT(*) AS total_count,
    COUNT(*) FILTER (
        WHERE completed_at IS NULL) AS incomplete_count,
    COUNT(*) FILTER (
        WHERE deadline >= week.week_start
          AND deadline < week.week_start + interval '7 days'
    ) AS due_this_week_count
FROM todo_items, week
WHERE student_id = $1
