SELECT id, subtotal, student_id, entry_type, essay_id, completed_at, refund, status, actor_id, created_at, updated_at
FROM public.essay_review_transaction
WHERE id = $1
FOR UPDATE