SELECT EXISTS (
    SELECT 1
    FROM public.media_access
    WHERE student_id = $1 AND media_id = $2
);
