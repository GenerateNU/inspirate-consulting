-- shift everything ranked below a removed application up by one so ranks stay consecutive
UPDATE public.personal_college_applications
SET rank = rank - 1
WHERE student_id = $1 AND rank > $2
