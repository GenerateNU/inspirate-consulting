-- $4 is an optional search term ('' returns everything) and $5 the minimum pg_trgm word_similarity:
-- how closely the search term matches the best-matching part of the title or description.
-- Results are ordered by best match; with no search term every score is 0, so they fall back to title.
-- due_date is the deadline of the student's todo that links this media (the todo that granted
-- access). If the media was assigned more than once, the most recently created todo wins.
SELECT m.id, m.title, m.description, m.length_in_mins, m.school_year, m.s3_key, m.youtube_video_id,
       ma.started_at, ma.completed_at,
       (
           SELECT t.deadline
           FROM public.todo_items t
           WHERE t.student_id = ma.student_id
             AND t.media_id = m.id
           ORDER BY t.created_at DESC, t.id DESC
           LIMIT 1
       ) AS due_date
FROM public.media m
JOIN public.media_access ma ON ma.media_id = m.id
WHERE ma.student_id = $1
  AND ($4::text = '' OR word_similarity($4, m.title) > $5 OR word_similarity($4, m.description) > $5)
ORDER BY GREATEST(word_similarity($4, m.title), word_similarity($4, m.description)) DESC, m.title, m.id
LIMIT $2 OFFSET $3;
