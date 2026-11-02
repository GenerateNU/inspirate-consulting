UPDATE public.student
SET
    year = COALESCE($2, year),
    organization = COALESCE($3, organization),
    gpa = COALESCE($4, gpa),
    student_interest = COALESCE($5, student_interest),
    updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, year, organization, gpa, review_balance, counselor_id, student_interest
