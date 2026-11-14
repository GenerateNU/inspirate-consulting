create table public.vector_relationship (
    vector_id bigint not null references public.vector_table(vector_id) on delete cascade,
    message_id uuid unique references public.chat_messages(id) on delete cascade,
    todo_id uuid unique references public.todo_items(id) on delete cascade,
    created_at timestamptz not null default now(),
    check (num_nonnulls(message_id, todo_id) = 1) -- exactly one of message_id or todo_id must be set
);

-- vector_id isn't unique (several messages/todos can share one vector)
-- index it for the on delete cascade from vector_table and the orphan lookup in the trigger below.
create index vector_relationship_vector_id_idx on public.vector_relationship (vector_id);

-- When a relationship row is removed (e.g. its message/todo was deleted and the FK cascaded),
-- delete the embedding once nothing else references it, so it doesn't linger in similarity searches.
create function public.delete_orphaned_vector()
returns trigger
language plpgsql
as $$
begin
    delete from public.vector_table vt
    where vt.vector_id = old.vector_id
      and not exists (
          select 1 from public.vector_relationship vr where vr.vector_id = old.vector_id
      );
    return old;
end;
$$;

create trigger vector_relationship_delete_vector
after delete on public.vector_relationship
for each row execute function public.delete_orphaned_vector();
