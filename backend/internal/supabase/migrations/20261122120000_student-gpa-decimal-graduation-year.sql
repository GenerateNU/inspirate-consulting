-- GPA was stored as a whole number, so values like 3.85 could not be saved.
alter table public.student
    alter column gpa type numeric(4, 2) using gpa::numeric(4, 2);

alter table public.student
    add column graduation_year smallint;

update public.student
set graduation_year = (
    extract(year from now())::int
    + case when extract(month from now()) >= 7 then 1 else 0 end
    + case year
        when 'senior' then 0
        when 'junior' then 1
        when 'sophomore' then 2
        when 'freshman' then 3
    end
);

alter table public.student
    alter column graduation_year set not null,
    drop column year;
