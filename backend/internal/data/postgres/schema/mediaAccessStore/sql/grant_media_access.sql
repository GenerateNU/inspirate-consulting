INSERT INTO public.media_access (
    student_id, media_id
) VALUES (
    $1, $2
)
RETURNING id, student_id, media_id