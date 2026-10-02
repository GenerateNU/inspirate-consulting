UPDATE public.essay_review_transaction
SET completed_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING id, subtotal, student_id, entry_type, essay_id, completed_at, refund, status, actor_id, created_at, updated_at
