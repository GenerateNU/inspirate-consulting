-- The caller's own students, deduplicated.
SELECT DISTINCT s.id
FROM public.student s
WHERE s.counselor_id = $1
  AND (
        $2::boolean
        OR s.id = ANY($3::uuid[])
        OR s.graduation_year = ANY($4::int[])
      )