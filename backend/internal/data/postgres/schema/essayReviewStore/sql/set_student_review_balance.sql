UPDATE public.student
SET review_balance = $2, updated_at = NOW()
WHERE id = $1