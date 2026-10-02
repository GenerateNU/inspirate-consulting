INSERT INTO public.essay_review_transaction (
    subtotal, student_id, entry_type
) VALUES (
    $1, $2, 'adjustment'
)