Insert into users(name, supabase_id, pfp_key, reset_time, email)
VALUES ($1, $2, $3,  $4, $5)
RETURNING id, name, supabase_id, pfp_key, reset_time, email