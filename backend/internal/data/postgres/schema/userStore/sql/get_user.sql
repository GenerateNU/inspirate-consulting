SELECT u.id, u.name, u.supabase_id, u.pfp_key, u.reset_time, u.email, u.expires_at
FROM users u
WHERE u.id = $1;
