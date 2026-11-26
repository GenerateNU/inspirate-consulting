-- ============================================
-- 12. CHAT MESSAGES
-- ============================================
-- sender_id / recipient_id are users.id. user_low / user_high are generated.
INSERT INTO chat_messages (id, sender_id, recipient_id, message, created_at, edited_at, read_at) VALUES
-- Maya Patel <-> Ethan Nguyen
('90000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000004', 'Hi Ethan! I left comments on your CMU "why us" draft. The robotics lab paragraph is strong, but tie it back to a specific CMU program.', '2026-11-20 10:15:00-05', NULL, '2026-11-20 16:02:00-05'),
('90000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000001', 'Thanks! Would mentioning the Robotics Institute summer program work?',                                                              '2026-11-20 16:05:00-05', NULL, '2026-11-20 17:30:00-05'),
('90000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000004', 'Perfect. Just name a faculty member or project you actually want to work with.',                                                    '2026-11-20 17:32:00-05', NULL, '2026-11-20 18:00:00-05'),
('90000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000001', 'Got it, I''ll have a revision up by Sunday.',                                                                                   '2026-11-20 18:01:00-05', '2026-11-20 18:03:00-05', NULL),
-- Maya Patel <-> Priya Shah
('90000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000001', 'Hi Ms. Patel, I archived the first BU supplement and started over. Can you look at the new version when you get a chance?',            '2026-11-16 15:25:00-05', NULL, '2026-11-16 19:00:00-05'),
('90000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000005', 'Of course. I refunded the review credit for the old draft, so you still have 3 left.',                                             '2026-11-16 19:05:00-05', NULL, '2026-11-17 07:45:00-05'),
-- Maya Patel <-> Lucas Bennett
('90000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000006', 'Lucas, please finish the "Building a Balanced College List" video before our meeting next week.',                                  '2026-11-18 09:00:00-05', NULL, NULL),
-- Daniel Okafor <-> Amara Johnson
('90000000-0000-0000-0000-000000000008', '10000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000002', 'Mr. Okafor, the UC app is due 11/30. Is my leadership PIQ ready to submit?',                                                    '2026-11-22 21:15:00-05', NULL, '2026-11-23 08:30:00-05'),
('90000000-0000-0000-0000-000000000009', '10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000007', 'Almost. Cut the intro by about 50 words and lead with the voter registration drive. I''ll do a final read Wednesday.',       '2026-11-23 08:35:00-05', NULL, NULL),
-- Daniel Okafor <-> Kenji Watanabe
('90000000-0000-0000-0000-00000000000a', '10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000008', 'Congrats on getting all your EA apps in! Next up is the CMU RD essay.',                                                       '2026-11-02 12:00:00-05', NULL, '2026-11-02 18:20:00-05'),
('90000000-0000-0000-0000-00000000000b', '10000000-0000-0000-0000-000000000008', '10000000-0000-0000-0000-000000000002', 'Thank you! I just submitted it for review.',                                                                                    '2026-11-22 14:35:00-05', NULL, NULL),
-- Sofia Ramirez <-> Noah Cohen
('90000000-0000-0000-0000-00000000000c', '10000000-0000-0000-0000-00000000000a', '10000000-0000-0000-0000-000000000003', 'Northeastern ED decisions come out mid-December, right?',                                                                     '2026-11-12 19:35:00-05', NULL, '2026-11-13 09:10:00-05'),
('90000000-0000-0000-0000-00000000000d', '10000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-00000000000a', 'Yes, usually mid-December. Keep working on your BU and Harvard RD essays just in case.',                                           '2026-11-13 09:12:00-05', NULL, '2026-11-13 12:00:00-05'),
-- Counselor <-> counselor
('90000000-0000-0000-0000-00000000000e', '10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000001', 'Do you have a good FAFSA walkthrough? A few of my seniors are asking.',                                                      '2026-11-19 11:00:00-05', NULL, '2026-11-19 11:20:00-05'),
('90000000-0000-0000-0000-00000000000f', '10000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000002', 'The "Financial Aid and the FAFSA" module covers it. I''m assigning it to Ethan and Priya this week.',                        '2026-11-19 11:22:00-05', NULL, '2026-11-19 11:40:00-05');
