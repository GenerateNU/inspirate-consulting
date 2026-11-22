-- $3 is an optional search term ('' returns everything) and $4 the minimum pg_trgm similarity.
-- Results are ordered by best match; with no search term every score is 0, so they fall back to title.
SELECT id, title, description, length_in_mins, school_year, s3_key, youtube_video_id
FROM public.media
WHERE $3::text = '' OR similarity(title, $3) > $4 OR similarity(description, $3) > $4
ORDER BY GREATEST(similarity(title, $3), similarity(description, $3)) DESC, title, id
LIMIT $1 OFFSET $2;
