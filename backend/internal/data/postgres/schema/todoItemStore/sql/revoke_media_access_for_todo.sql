DELETE FROM public.media_access
WHERE student_id = $1
  AND media_id = $2
  AND NOT EXISTS (
      SELECT 1
      FROM public.todo_items
      WHERE student_id = $1 AND media_id = $2
  );
