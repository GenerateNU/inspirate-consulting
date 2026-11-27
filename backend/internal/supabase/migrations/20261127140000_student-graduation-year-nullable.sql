-- Students only need a counselor at sign up; they add their graduation year later.
alter table public.student
    alter column graduation_year drop not null;
