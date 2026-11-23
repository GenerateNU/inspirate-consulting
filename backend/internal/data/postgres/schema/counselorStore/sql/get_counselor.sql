SELECT c.id, c.user_id, c.custom_link, u.expires_at
FROM public.counselor AS c
JOIN public.users AS u ON u.id = c.user_id
WHERE c.id = $1
