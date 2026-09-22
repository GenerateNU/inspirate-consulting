create table essays (
    id uuid primary key default gen_random_uuid(),
    student_id uuid not null references students(id), 
    type text not null,
    college_id uuid references global_college_info(id),
    link_to_content text not null
    status text not null default 'Draft'
)