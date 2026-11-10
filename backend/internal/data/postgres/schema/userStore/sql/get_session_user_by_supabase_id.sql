SELECT u.id, u.name, u.supabase_id, u.pfp_key, u.reset_time, s.id, c.id
FROM users u
LEFT JOIN student s ON s.user_id = u.id
LEFT JOIN counselor c ON c.user_id = u.id
WHERE u.supabase_id = $1;
