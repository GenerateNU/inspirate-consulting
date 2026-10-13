Insert into users(name, supabase_id, pfp_key, reset_time)
VALUES ($1, $2, $3,  $4)
RETURNING id, name, supabase_id, pfp_key, reset_time