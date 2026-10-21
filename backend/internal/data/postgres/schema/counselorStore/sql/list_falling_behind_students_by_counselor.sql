SELECT DISTINCT todo.student_id
FROM public.todo_items AS todo
JOIN public.student AS student ON student.id = todo.student_id
WHERE student.counselor_id = $1
  AND todo.completed_at IS NULL
  AND todo.deadline < now()
