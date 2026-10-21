SELECT EXISTS (
    SELECT 1
    FROM student
    WHERE id = $1 AND counselor_id = $2
)
