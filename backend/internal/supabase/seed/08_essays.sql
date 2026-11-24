-- ============================================
-- 8. ESSAYS
-- ============================================
INSERT INTO essays (id, student_id, type, college_id, link_to_content, status, essay_group_id) VALUES
-- Ethan Nguyen
('50000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'personal-statement', NULL, 'https://docs.google.com/document/d/seed-ethan-personal-statement',    'Submitted', '40000000-0000-0000-0000-000000000001'),
('50000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000001', 'short-answer',       2,    'https://docs.google.com/document/d/seed-ethan-mit-community',         'Submitted', '40000000-0000-0000-0000-000000000002'),
('50000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000001', 'why-us',             10,   'https://docs.google.com/document/d/seed-ethan-cmu-why-us',            'Review',    NULL),
('50000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000001', 'supplemental',       9,    'https://docs.google.com/document/d/seed-ethan-gt-supplement',         'Draft',     NULL),
-- Priya Shah
('50000000-0000-0000-0000-000000000005', '30000000-0000-0000-0000-000000000002', 'personal-statement', NULL, 'https://docs.google.com/document/d/seed-priya-personal-statement',    'Submitted', '40000000-0000-0000-0000-000000000003'),
('50000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000002', 'supplemental',       3,    'https://docs.google.com/document/d/seed-priya-harvard-intellectual',  'Submitted', '40000000-0000-0000-0000-000000000004'),
('50000000-0000-0000-0000-000000000007', '30000000-0000-0000-0000-000000000002', 'why-us',             5,    'https://docs.google.com/document/d/seed-priya-tufts-why-us',          'Review',    NULL),
('50000000-0000-0000-0000-000000000008', '30000000-0000-0000-0000-000000000002', 'supplemental',       4,    'https://docs.google.com/document/d/seed-priya-bu-supplement-v1',      'Archived',  NULL),
('50000000-0000-0000-0000-000000000009', '30000000-0000-0000-0000-000000000002', 'supplemental',       4,    'https://docs.google.com/document/d/seed-priya-bu-supplement-v2',      'Draft',     NULL),
-- Lucas Bennett
('50000000-0000-0000-0000-00000000000a', '30000000-0000-0000-0000-000000000003', 'personal-statement', NULL, 'https://docs.google.com/document/d/seed-lucas-brainstorm',           'Draft',     NULL),
-- Amara Johnson
('50000000-0000-0000-0000-00000000000b', '30000000-0000-0000-0000-000000000004', 'personal-statement', NULL, 'https://docs.google.com/document/d/seed-amara-personal-statement',    'Submitted', NULL),
('50000000-0000-0000-0000-00000000000c', '30000000-0000-0000-0000-000000000004', 'uc-piq',             8,    'https://docs.google.com/document/d/seed-amara-piq-leadership',        'Review',    '40000000-0000-0000-0000-000000000005'),
('50000000-0000-0000-0000-00000000000d', '30000000-0000-0000-0000-000000000004', 'uc-piq',             8,    'https://docs.google.com/document/d/seed-amara-piq-community',         'Draft',     '40000000-0000-0000-0000-000000000005'),
('50000000-0000-0000-0000-00000000000e', '30000000-0000-0000-0000-000000000004', 'short-answer',       7,    'https://docs.google.com/document/d/seed-amara-stanford-roommate',     'Submitted', NULL),
-- Kenji Watanabe
('50000000-0000-0000-0000-00000000000f', '30000000-0000-0000-0000-000000000005', 'personal-statement', NULL, 'https://docs.google.com/document/d/seed-kenji-personal-statement',    'Submitted', '40000000-0000-0000-0000-000000000006'),
('50000000-0000-0000-0000-000000000010', '30000000-0000-0000-0000-000000000005', 'why-major',          13,   'https://docs.google.com/document/d/seed-kenji-purdue-why-major',      'Submitted', '40000000-0000-0000-0000-000000000006'),
('50000000-0000-0000-0000-000000000011', '30000000-0000-0000-0000-000000000005', 'why-us',             10,   'https://docs.google.com/document/d/seed-kenji-cmu-why-us',            'Review',    NULL),
-- Noah Cohen
('50000000-0000-0000-0000-000000000012', '30000000-0000-0000-0000-000000000007', 'personal-statement', NULL, 'https://docs.google.com/document/d/seed-noah-personal-statement',     'Submitted', '40000000-0000-0000-0000-000000000007'),
('50000000-0000-0000-0000-000000000013', '30000000-0000-0000-0000-000000000007', 'why-us',             1,    'https://docs.google.com/document/d/seed-noah-neu-why-us',             'Submitted', '40000000-0000-0000-0000-000000000007'),
('50000000-0000-0000-0000-000000000014', '30000000-0000-0000-0000-000000000007', 'supplemental',       4,    'https://docs.google.com/document/d/seed-noah-bu-supplement',          'Draft',     NULL);
