SELECT student.id, users.name, student.year,

        (SELECT count(*)::int
        FROM public.essay_review_transaction AS review
        WHERE review.student_id = student.id
        AND review.status = 'open') AS pending_review_count,

       -- Incomplete to-do items
       (SELECT count(*)::int
        FROM public.todo_items AS todo
        WHERE todo.student_id = student.id
        AND todo.completed_at IS NULL) AS open_task_count,

       -- Messages from the student to this counselor that haven't been read
       (SELECT count(*)::int
       FROM public.chat_messages AS message
       WHERE message.sender_id = student.user_id
       AND message.recipient_id = (SELECT counselor.user_id
       FROM public.counselor AS counselor
       WHERE counselor.id = $1)
       AND message.read_at IS NULL) AS unread_message_count

FROM public.student AS student
JOIN public.users AS users ON users.id = student.user_id
WHERE student.counselor_id = $1
