INSERT INTO public.essay_review_transaction (
    subtotal, student_id, entry_type, essay_id
) VALUES (
    $1, $2, 'refund', $3
)
RETURNING id, subtotal, student_id, entry_type, essay_id, completed_at, refund, status, actor_id, created_at, updated_at