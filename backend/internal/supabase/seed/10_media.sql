-- ============================================
-- 10. MEDIA
-- ============================================
-- Exactly one of s3_key / youtube_video_id is set per row.
-- The S3 keys and YouTube ids are placeholders and won't resolve to real videos.
INSERT INTO media (id, title, description, length_in_mins, school_year, s3_key, youtube_video_id) VALUES
('70000000-0000-0000-0000-000000000001', 'Welcome to Inspirate',                       'How the platform works: tasks, essays, college list, and messaging your counselor.',        6,    NULL,        'video-bucket/welcome-to-inspirate.mp4',     NULL),
('70000000-0000-0000-0000-000000000002', 'Writing Your Personal Statement',            'Brainstorming, structure, and revision strategies for the Common App essay.',               24,   'senior',    'video-bucket/personal-statement.mp4',       NULL),
('70000000-0000-0000-0000-000000000003', 'Building a Balanced College List',           'How to pick reach, target, and safety schools that fit academically and financially.',     18,   'junior',    'video-bucket/balanced-college-list.mp4',    NULL),
('70000000-0000-0000-0000-000000000004', 'Early Decision vs. Early Action Explained',  'ED, ED2, EA, REA and rolling admissions, and when each one makes sense.',                   NULL, 'junior',    NULL,                                        'seedEDvsEA01'),
('70000000-0000-0000-0000-000000000005', 'Choosing Extracurriculars That Matter',      'Depth over breadth: committing to activities and growing into leadership.',                 15,   'sophomore', 'video-bucket/extracurriculars.mp4',         NULL),
('70000000-0000-0000-0000-000000000006', 'Financial Aid and the FAFSA',                'Walkthrough of the FAFSA and CSS Profile, plus reading an aid award letter.',                NULL, 'senior',    NULL,                                        'seedFAFSA001'),
('70000000-0000-0000-0000-000000000007', 'Preparing for College Interviews',           'Common alumni interview questions and how to talk about yourself with confidence.',         12,   'senior',    'video-bucket/college-interviews.mp4',       NULL),
('70000000-0000-0000-0000-000000000008', 'SAT vs. ACT: Which Test Is Right for You?',  'Comparing the two tests, test-optional policies, and building a prep timeline.',             20,   'sophomore', NULL,                                        'seedSATACT01'),
('70000000-0000-0000-0000-000000000009', 'Making the Most of Freshman Year',           'Building study habits, exploring clubs, and starting a brag sheet early.',                  10,   'freshman',  'video-bucket/freshman-year.mp4',            NULL);
