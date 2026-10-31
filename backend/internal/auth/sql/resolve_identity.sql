SELECT u.id, s.id, c.id, u.expires_at
FROM users u
LEFT JOIN student s ON s.user_id = u.id
LEFT JOIN counselor c ON c.user_id = u.id
WHERE u.supabase_id = $1
