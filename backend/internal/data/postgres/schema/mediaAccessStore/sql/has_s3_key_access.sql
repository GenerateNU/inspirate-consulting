SELECT EXISTS (
    SELECT 1
    FROM public.media_access ma
    JOIN public.media m ON m.id = ma.media_id
    WHERE ma.student_id = $1 AND m.s3_key = $2
);
