UPDATE public.student
SET
    year = COALESCE($2, year),
    organization = COALESCE($3, organization),
    gpa = COALESCE($4, gpa),
    student_interest = COALESCE($5, student_interest),
    high_school = COALESCE($6, high_school),
    standardized_test_score = COALESCE($7, standardized_test_score),
    intended_major = COALESCE($8, intended_major),
    intended_minor = COALESCE($9, intended_minor),
    geographic_preferences = COALESCE($10, geographic_preferences),
    updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, year, organization, gpa, review_balance, counselor_id, student_interest,
    high_school, standardized_test_score, intended_major, intended_minor, geographic_preferences
