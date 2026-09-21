SELECT u.name, u.supabase_id, u.pfp_key
FROM users u
WHERE u.id = $1;