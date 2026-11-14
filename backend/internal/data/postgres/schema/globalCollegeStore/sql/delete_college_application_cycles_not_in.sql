DELETE FROM public.college_application_cycles
WHERE global_college_id = $1 AND NOT (cycle = ANY($2))
