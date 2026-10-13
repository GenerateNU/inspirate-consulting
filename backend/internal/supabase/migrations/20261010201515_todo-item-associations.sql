create extension if not exists pg_trgm;

alter table public.todo_items
    add column essay_id uuid references public.essays(id) on delete set null,
    add column media_id uuid references public.media(id) on delete set null,
    add column global_college_id bigint references public.global_colleges(id) on delete set null,
    add constraint todo_items_single_association
        check (num_nonnulls(essay_id, media_id, global_college_id) <= 1);

create index todo_items_student_id_idx on public.todo_items (student_id);
create index todo_items_essay_id_idx on public.todo_items (essay_id) where essay_id is not null;
create index todo_items_media_id_idx on public.todo_items (media_id) where media_id is not null;
create index todo_items_global_college_id_idx on public.todo_items (global_college_id) where global_college_id is not null;
-- gin index to support better performance on text search on todo_description column
create index todo_items_description_trgm_idx on public.todo_items using gin (todo_description gin_trgm_ops);
