SELECT c.id, c.user_id, u.name, u.pfp_key
FROM public.counselor AS c
JOIN public.users AS u
ON c.user_id = u.id
WHERE ($1::text is NULL or u.name ILIKE '%' || $1 || '%')
AND (u.name, c.id) > ($2, $3)
ORDER BY u.name, c.id DESC
LIMIT $4