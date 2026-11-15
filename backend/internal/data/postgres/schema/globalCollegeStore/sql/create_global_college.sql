INSERT INTO public.global_colleges (
    school_name, school_location
) VALUES (
    $1, $2
)
RETURNING id, created_at, updated_at, school_name, school_location
