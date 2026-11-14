UPDATE users
SET
    name = COALESCE($2, name),
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, supabase_id, pfp_key, reset_time, email, expires_at
