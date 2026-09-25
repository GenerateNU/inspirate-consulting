Insert into users(name, supabase_id, pfp_key)
VALUES ($1, $2, $3)
RETURNING id, name, supabase_id, pfp_key