SELECT id, user_id, year, organization, gpa, review_balance, counselor_id
FROM public.student
WHERE id = $1
FOR UPDATE