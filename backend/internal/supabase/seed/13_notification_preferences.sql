-- ============================================
-- 13. NOTIFICATION PREFERENCES
-- ============================================
INSERT INTO notification_preferences (user_id, email_enabled, weekly_summary_enabled, due_date_notifications_enabled, days_before_due, notify_past_due) VALUES
-- Counselors
('10000000-0000-0000-0000-000000000001', true,  true,  false, 1, true),
('10000000-0000-0000-0000-000000000002', true,  true,  false, 1, false),
('10000000-0000-0000-0000-000000000003', true,  false, false, 1, false),
-- Students
('10000000-0000-0000-0000-000000000004', true,  true,  true,  2, true),
('10000000-0000-0000-0000-000000000005', true,  true,  true,  3, true),
('10000000-0000-0000-0000-000000000006', true,  false, true,  1, true),
('10000000-0000-0000-0000-000000000007', true,  true,  true,  1, true),
('10000000-0000-0000-0000-000000000008', true,  true,  true,  7, false),
('10000000-0000-0000-0000-000000000009', false, false, false, 1, false),
('10000000-0000-0000-0000-00000000000a', true,  true,  true,  0, true),
('10000000-0000-0000-0000-00000000000b', true,  true,  true,  1, true);
