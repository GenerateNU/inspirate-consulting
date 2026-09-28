create table public.media (
    id uuid primary key default gen_random_uuid(),
    title text not null,
    description text not null,
    length_in_mins integer not null,
    school_year integer,
    s3_key text not null

);

create table public.media_access (
    id uuid primary key default gen_random_uuid(),
    student_id uuid not null,
    media_id uuid not null references public.media(id)
);