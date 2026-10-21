SELECT review.id, review.essay_id, essay.type, essay.link_to_content, student.id, users.name, review.created_at
FROM public.essay_review_transaction AS review
JOIN public.student AS student ON student.id = review.student_id
JOIN public.essays AS essay ON essay.id = review.essay_id
JOIN public.users AS users ON users.id = student.user_id
WHERE student.counselor_id = $1
  AND review.status = 'open'
ORDER BY review.created_at ASC, review.id ASC
