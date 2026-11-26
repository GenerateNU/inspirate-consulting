-- NULL params leave the column unchanged. $4 (s3_key) and $5 (youtube_video_id) are mutually
-- exclusive: setting one clears the other so the media_single_source constraint still holds.
UPDATE public.media
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    s3_key = CASE
        WHEN $4::text IS NOT NULL THEN $4
        WHEN $5::text IS NOT NULL THEN NULL
        ELSE s3_key
    END,
    youtube_video_id = CASE
        WHEN $5::text IS NOT NULL THEN $5
        WHEN $4::text IS NOT NULL THEN NULL
        ELSE youtube_video_id
    END
WHERE id = $1
RETURNING id, title, description, length_in_mins, school_year, s3_key, youtube_video_id;
