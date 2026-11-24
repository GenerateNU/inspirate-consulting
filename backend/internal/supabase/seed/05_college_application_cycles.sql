-- ============================================
-- 5. COLLEGE APPLICATION CYCLES
-- ============================================
-- Fall 2027 admissions cycle. Only ROLLING rows use rolling_start/rolling_end instead of a deadline.
INSERT INTO college_application_cycles (global_college_id, cycle, deadline, rolling_start, rolling_end) VALUES
-- Northeastern
(1,  'ED',      '2026-11-01 23:59:00-05', NULL,         NULL),
(1,  'EA',      '2026-11-01 23:59:00-05', NULL,         NULL),
(1,  'ED2',     '2027-01-01 23:59:00-05', NULL,         NULL),
(1,  'RD',      '2027-01-01 23:59:00-05', NULL,         NULL),
-- MIT
(2,  'EA',      '2026-11-01 23:59:00-05', NULL,         NULL),
(2,  'RD',      '2027-01-05 23:59:00-05', NULL,         NULL),
-- Harvard
(3,  'REA',     '2026-11-01 23:59:00-05', NULL,         NULL),
(3,  'RD',      '2027-01-01 23:59:00-05', NULL,         NULL),
-- Boston University
(4,  'ED',      '2026-11-01 23:59:00-05', NULL,         NULL),
(4,  'ED2',     '2027-01-04 23:59:00-05', NULL,         NULL),
(4,  'RD',      '2027-01-04 23:59:00-05', NULL,         NULL),
-- Tufts
(5,  'ED',      '2026-11-01 23:59:00-05', NULL,         NULL),
(5,  'ED2',     '2027-01-04 23:59:00-05', NULL,         NULL),
(5,  'RD',      '2027-01-04 23:59:00-05', NULL,         NULL),
-- Michigan
(6,  'EA',      '2026-11-01 23:59:00-05', NULL,         NULL),
(6,  'RD',      '2027-02-01 23:59:00-05', NULL,         NULL),
-- Stanford
(7,  'REA',     '2026-11-01 23:59:00-08', NULL,         NULL),
(7,  'RD',      '2027-01-05 23:59:00-08', NULL,         NULL),
-- UC Berkeley
(8,  'RD',      '2026-11-30 23:59:00-08', NULL,         NULL),
-- Georgia Tech
(9,  'EA',      '2026-11-01 23:59:00-05', NULL,         NULL),
(9,  'RD',      '2027-01-05 23:59:00-05', NULL,         NULL),
-- Carnegie Mellon
(10, 'ED',      '2026-11-01 23:59:00-05', NULL,         NULL),
(10, 'ED2',     '2027-01-03 23:59:00-05', NULL,         NULL),
(10, 'RD',      '2027-01-03 23:59:00-05', NULL,         NULL),
-- UMass Amherst
(11, 'EA',      '2026-11-05 23:59:00-05', NULL,         NULL),
(11, 'RD',      '2027-01-15 23:59:00-05', NULL,         NULL),
-- Arizona State
(12, 'ROLLING', NULL,                     '2026-08-01', '2027-05-01'),
-- Purdue
(13, 'EA',      '2026-11-01 23:59:00-05', NULL,         NULL),
(13, 'RD',      '2027-01-15 23:59:00-05', NULL,         NULL);
