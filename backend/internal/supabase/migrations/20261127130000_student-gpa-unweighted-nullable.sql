
alter table public.student
    drop column gpa;

alter table public.student
    add column gpa double precision
        constraint student_gpa_unweighted_range check (gpa between 0 and 4);
