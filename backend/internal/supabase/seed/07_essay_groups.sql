-- ============================================
-- 7. ESSAY GROUPS
-- ============================================
-- name is globally unique, so each group name carries the student's name.
INSERT INTO essay_groups (id, name, description, student_id) VALUES
('40000000-0000-0000-0000-000000000001', 'Ethan N. - Common App',                'Personal statement drafts for the Common App',             '30000000-0000-0000-0000-000000000001'),
('40000000-0000-0000-0000-000000000002', 'Ethan N. - MIT Short Responses',       'MIT short-answer prompts (community, challenge, fun)',     '30000000-0000-0000-0000-000000000001'),
('40000000-0000-0000-0000-000000000003', 'Priya S. - Common App',                'Personal statement and activity descriptions',             '30000000-0000-0000-0000-000000000002'),
('40000000-0000-0000-0000-000000000004', 'Priya S. - Harvard Supplements',       'Harvard REA supplemental questions',                       '30000000-0000-0000-0000-000000000002'),
('40000000-0000-0000-0000-000000000005', 'Amara J. - UC Personal Insight Questions', 'Four PIQs for the UC application',                     '30000000-0000-0000-0000-000000000004'),
('40000000-0000-0000-0000-000000000006', 'Kenji W. - Common App',                'Personal statement and engineering "why major" essays',   '30000000-0000-0000-0000-000000000005'),
('40000000-0000-0000-0000-000000000007', 'Noah C. - Northeastern ED',            NULL,                                                       '30000000-0000-0000-0000-000000000007');
