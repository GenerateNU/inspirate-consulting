Insert into users(name, supabase_id, pfp_key, reset_time, email, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, supabase_id, pfp_key, reset_time, email, expires_at
