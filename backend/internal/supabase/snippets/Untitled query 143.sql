-- Users: one for the student, one for their counselor
insert into public.users (id, name, email) values
  ('00000000-0000-0000-0000-000000000001', 'Risha Jones', 'risha@test.com'),
  ('00000000-0000-0000-0000-000000000003', 'Test Counselor', 'counselor@test.com');

-- Counselor
insert into public.counselor (user_id)
values ('00000000-0000-0000-0000-000000000003');

-- Student, using the same ID as the auth stub
insert into public.student (id, user_id, year, gpa, review_balance, counselor_id, organization, student_interest)
values (
  '00000000-0000-0000-0000-000000000002',
  '00000000-0000-0000-0000-000000000001',
  'senior', 4, 0,
  (select id from public.counselor where user_id = '00000000-0000-0000-0000-000000000003'),
  'Jericho High School',
  'Data Science'
);

-- Todo items: some completed, some not
insert into public.todo_items (student_id, user_id, todo_description, completed_at) values
  ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'Write Georgetown supplemental', null),
  ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'Revise Parsons edits', null),
  ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'Common App draft', now()),
  ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 'Upload letter of rec', now());