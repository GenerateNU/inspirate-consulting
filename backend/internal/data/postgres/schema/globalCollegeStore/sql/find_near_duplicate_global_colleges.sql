SELECT id, created_at, updated_at, school_name, school_location
FROM public.global_colleges
WHERE similarity(school_name, $1) > $3 AND similarity(school_location, $2) > $3
ORDER BY similarity(school_name, $1) DESC