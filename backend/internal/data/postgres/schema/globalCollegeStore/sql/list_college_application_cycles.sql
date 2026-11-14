SELECT global_college_id, cycle, deadline, rolling_start::text, rolling_end::text
FROM public.college_application_cycles
WHERE global_college_id = ANY($1)
ORDER BY global_college_id, array_position(ARRAY['ED', 'ED2', 'REA', 'EA', 'RD', 'ROLLING'], cycle)
