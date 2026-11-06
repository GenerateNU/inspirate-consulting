UPDATE public.media
SET title = $2
WHERE id = $1
RETURNING id, title, description, length_in_mins, school_year, s3_key;
