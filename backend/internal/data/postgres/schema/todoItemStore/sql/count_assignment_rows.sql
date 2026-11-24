-- Does this assignment exist, and did this counselor create it?
SELECT count(*)
FROM public.todo_items
WHERE assignment_id = $1 AND user_id = $2