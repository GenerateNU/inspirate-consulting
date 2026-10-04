UPDATE users
SET reset_time = $2
WHERE id = $1;
