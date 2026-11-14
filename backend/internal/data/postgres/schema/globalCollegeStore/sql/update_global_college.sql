UPDATE public.global_colleges
SET
    school_name = $2,
    school_location = $3,
    updated_at = now()
WHERE id = $1
RETURNING id, created_at, updated_at, school_name, school_location
