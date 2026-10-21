SELECT id, created_at, updated_at, school_name, school_location, ea_deadline, ed_deadline, rd_deadline
FROM public.global_colleges
WHERE similarity(school_name, $1) > $2 OR similarity(school_location, $1) > $2
ORDER BY GREATEST(similarity(school_name, $1), similarity(school_location, $1)) DESC
LIMIT $3;