-- ============================================
-- 4. GLOBAL COLLEGES
-- ============================================
INSERT INTO global_colleges (id, school_name, school_location, created_at, updated_at) VALUES
(1,  'Northeastern University',                'Boston, MA',          '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(2,  'Massachusetts Institute of Technology',  'Cambridge, MA',       '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(3,  'Harvard University',                     'Cambridge, MA',       '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(4,  'Boston University',                      'Boston, MA',          '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(5,  'Tufts University',                       'Medford, MA',         '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(6,  'University of Michigan',                 'Ann Arbor, MI',       '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(7,  'Stanford University',                    'Stanford, CA',        '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(8,  'University of California, Berkeley',     'Berkeley, CA',        '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(9,  'Georgia Institute of Technology',        'Atlanta, GA',         '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(10, 'Carnegie Mellon University',             'Pittsburgh, PA',      '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(11, 'University of Massachusetts Amherst',    'Amherst, MA',         '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(12, 'Arizona State University',               'Tempe, AZ',           '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04'),
(13, 'Purdue University',                      'West Lafayette, IN',  '2026-08-01 10:00:00-04', '2026-08-01 10:00:00-04');

-- explicit ids were inserted, so move the identity sequence past them
SELECT setval(pg_get_serial_sequence('public.global_colleges', 'id'), (SELECT MAX(id) FROM global_colleges));
