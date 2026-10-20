SELECT s.id, s.user_id, u.name, u.pfp_key, s.year, s.organization, s.gpa, s.review_balance, s.counselor_id
FROM public.student AS s
JOIN public.users AS u
ON s.user_id = u.id
WHERE ($1::text IS NULL OR u.name ILIKE '%' || $1 || '%')
AND ($2::text is NULL OR (u.name, s.id) > ($2, $3))
ORDER BY u.name, s.id
LIMIT $4
