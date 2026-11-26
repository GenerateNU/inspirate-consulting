SELECT id, user_id, graduation_year, organization, gpa, review_balance, counselor_id, student_interest,
    high_school, standardized_test_score, intended_major, intended_minor, geographic_preferences
FROM public.student
WHERE id = $1
