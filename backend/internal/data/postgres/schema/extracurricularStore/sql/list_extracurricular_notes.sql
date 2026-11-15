SELECT
    id,
    extracurricular_id,
    body,
    created_at,
    modified_at
FROM public.extracurricular_notes
WHERE extracurricular_id = $1
ORDER BY created_at DESC, id DESC;
