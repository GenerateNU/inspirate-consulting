SELECT activity.id, activity.student_id, activity.action, activity.essay_id, activity.essay_status,
    activity.global_college_id, activity.extracurricular_id, activity.todo_item_id, activity.created_at
FROM public.activity
JOIN public.student ON student.id = activity.student_id
WHERE student.counselor_id = $1
    AND ($2::uuid IS NULL OR activity.student_id = $2)
ORDER BY activity.created_at DESC
