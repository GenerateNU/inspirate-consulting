SELECT u.id, u.name, u.supabase_id, u.pfp_key, u.reset_time
FROM users u
WHERE u.supabase_id = $1;
