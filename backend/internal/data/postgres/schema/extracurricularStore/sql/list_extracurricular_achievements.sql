SELECT
    id,
    extracurricular_id,
    entry_date,
    body,
    created_at,
    modified_at
FROM public.extracurricular_achievements
WHERE extracurricular_id = $1
ORDER BY entry_date DESC, id DESC;
