INSERT INTO public.college_application_cycles (
    global_college_id, cycle, deadline, rolling_start, rolling_end
) VALUES (
    $1, $2, $3, $4::date, $5::date
)
ON CONFLICT (global_college_id, cycle) DO UPDATE
SET
    deadline = EXCLUDED.deadline,
    rolling_start = EXCLUDED.rolling_start,
    rolling_end = EXCLUDED.rolling_end,
    updated_at = now()
