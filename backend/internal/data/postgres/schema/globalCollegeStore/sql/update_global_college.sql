UPDATE public.global_colleges
SET
    school_name = $2,
    school_location = $3,
    ea_deadline = $4,
    ed_deadline = $5,
    rd_deadline = $6,
    updated_at = now()
WHERE id = $1
RETURNING id, school_name, school_location, ea_deadline, ed_deadline, rd_deadline, created_at, updated_at