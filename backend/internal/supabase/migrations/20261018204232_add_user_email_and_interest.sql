alter table public.users
    add column email TEXT NOT NULL DEFAULT 'useremail@provider.com';

alter table public.users
    alter column email drop default;

alter table public.student
    add column student_interest TEXT;