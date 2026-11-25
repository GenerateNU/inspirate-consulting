SELECT id, title, description, length_in_mins, school_year, s3_key, youtube_video_id
FROM public.media
WHERE id = $1
