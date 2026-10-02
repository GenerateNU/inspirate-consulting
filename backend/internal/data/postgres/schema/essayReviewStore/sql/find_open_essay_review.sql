SELECT id
FROM public.essay_review_transaction
WHERE essay_id = $1 AND status = 'open'
LIMIT 1