SELECT student.id, users.name, student.year
FROM public.student AS student
JOIN public.users AS users ON users.id = student.user_id
WHERE student.counselor_id = $1
