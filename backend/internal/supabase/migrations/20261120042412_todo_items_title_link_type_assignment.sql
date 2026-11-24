-- title: add nullable and backfill from todo_description
alter table public.todo_items add column title text;

update public.todo_items
set title = coalesce(nullif(left(trim(todo_description), 80), ''), 'Unititled Task');

alter table public.todo_items alter column title set not null;

-- extracurricular relation
alter table public.todo_items
    add column extracurricular_id bigint references public.extracurriculars(id) on delete set null;

create index todo_items_extracurriculars_id_idx 
    on public.todo_items (extracurricular_id) 
    where extracurricular_id is not null;

-- link type: nullable and backfilled based from which id is set 
create type todo_link_type as enum ('college', 'essay', 'extracurricular', 'module');

alter table public.todo_items add column link_type public.todo_link_type;

update public.todo_items
set link_type = case 
    when essay_id is not null then 'essay'::public.todo_link_type
    when media_id is not null then 'module'::public.todo_link_type
    when global_college_id is not null then 'college'::public.todo_link_type
end;

-- single association now includes extracurricular_id
alter table public.todo_items drop constraint todo_items_single_association;

alter table public.todo_items add constraint todo_items_single_association
    check (num_nonnulls(essay_id, media_id, global_college_id, extracurricular_id) <= 1);
    
-- check that link type matches if id is set
-- id going null (eg essay being deleted) still passes so link_type stays
-- task still shows
alter table public.todo_items add constraint todo_link_type_matches
    check (
        (essay_id is null or link_type = 'essay')
        and (media_id is null or link_type = 'module')
        and (extracurricular_id is null or link_type = 'extracurricular')
        and (global_college_id is null or link_type = 'college')
    );

-- assignment_id: shared by every row in one multi-student assign
alter table public.todo_items add column assignment_id uuid;

create index idx_todo_items_assignment_id
    on public.todo_items (assignment_id)
    where assignment_id is not null;    




