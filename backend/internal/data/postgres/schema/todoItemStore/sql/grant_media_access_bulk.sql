-- Give these students access to a video, skipping any who already have it.
INSERT INTO public.media_access (media_id, student_id)
SELECT $1::uuid, sid
FROM unnest($2::uuid[]) AS sid
WHERE NOT EXISTS (
    SELECT 1 FROM public.media_access ma
    WHERE ma.media_id = $1::uuid AND ma.student_id = sid
)