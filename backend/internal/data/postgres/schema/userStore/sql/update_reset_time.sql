UPDATE users
SET reset_time = $2
WHERE supabase_id = $1;
