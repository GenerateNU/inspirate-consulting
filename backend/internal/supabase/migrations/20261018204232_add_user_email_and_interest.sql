alter table public.users
    add column email TEXT NOT NULL DEFAULT 'useremail@placeholder.invalid';

alter table public.users
    alter column email drop default;

alter table public.student
    add column student_interest TEXT;

create unique index idx_user_email 
    on public.users (email)
    where email <> 'useremail@placeholder.invalid';