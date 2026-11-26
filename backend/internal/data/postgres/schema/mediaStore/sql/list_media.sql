-- $3 is an optional search term ('' returns everything) and $4 the minimum pg_trgm word_similarity:
-- how closely the search term matches the best-matching part of the title or description.
-- Results are ordered by best match; with no search term every score is 0, so they fall back to title.
SELECT id, title, description, length_in_mins, school_year, s3_key, youtube_video_id
FROM public.media
WHERE $3::text = '' OR word_similarity($3, title) > $4 OR word_similarity($3, description) > $4
ORDER BY GREATEST(word_similarity($3, title), word_similarity($3, description)) DESC, title, id
LIMIT $1 OFFSET $2;
