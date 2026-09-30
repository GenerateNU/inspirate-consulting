SELECT id AS transaction_id, essay_id, student_id, status, created_at AS requested_at, completed_at
FROM public.essay_review_transaction
WHERE essay_id = $1 AND entry_type = 'spend'
ORDER BY created_at DESC, id DESC
LIMIT 1