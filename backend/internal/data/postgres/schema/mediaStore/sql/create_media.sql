INSERT INTO public.media (
    title, description, length_in_mins, school_year, s3_key, youtube_video_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, title, description, length_in_mins, school_year, s3_key, youtube_video_id
