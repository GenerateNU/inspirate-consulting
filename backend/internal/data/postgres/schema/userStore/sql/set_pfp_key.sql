UPDATE users u
SET
    pfp_key = $2,
    updated_at = NOW()
FROM (SELECT id, pfp_key FROM users WHERE id = $1 FOR UPDATE) previous
WHERE u.id = previous.id
RETURNING u.id, u.name, u.supabase_id, u.pfp_key, u.reset_time, u.email, u.expires_at, previous.pfp_key
