INSERT INTO public.global_colleges (
    school_name, school_location, ea_deadline, ed_deadline, rd_deadline
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, created_at, updated_at, school_name, school_location, ea_deadline, ed_deadline, rd_deadline
