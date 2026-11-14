-- essays with review history are kept (the credit ledger points at them) but archived.
-- skipped if the student still has another application to the same college.
UPDATE public.essays e
SET status = 'Archived'
WHERE e.student_id = $1
  AND e.college_id = $2
  AND EXISTS (SELECT 1 FROM public.essay_review_transaction t WHERE t.essay_id = e.id)
  AND NOT EXISTS (
      SELECT 1 FROM public.personal_college_applications p
      WHERE p.student_id = $1 AND p.global_college_id = $2
  )
