SELECT id, title, description, length_in_mins, school_year, s3_key
FROM public.media
ORDER BY title, id
LIMIT $1 OFFSET $2;