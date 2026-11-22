UPDATE public.student
SET counselor_id = $2,
    updated_at = NOW()
WHERE user_id = $1
RETURNING id, user_id, graduation_year, organization, gpa, review_balance, counselor_id
