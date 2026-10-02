create table public.todo_items(
    id uuid primary key default gen_random_uuid(),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    student_id uuid not null,
    user_id uuid not null,
    todo_description text not null,
    completed_at timestamptz,
    deadline timestamptz
);