SELECT id, created_at, updated_at, school_name, school_location
FROM public.global_colleges
WHERE id = $1