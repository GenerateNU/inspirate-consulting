-- Embeddings for chat messages and todo items

-- pg_net sends the async HTTP request.
-- It only exists on Supabase, so skip it on plain Postgres.
do $$
begin
    if exists (select 1 from pg_available_extensions where name = 'pg_net') then
        create extension if not exists pg_net;
    end if;
end;
$$;

-- Replaces the vector for one message or todo item (or creates it) in one transaction.
-- Removing the old relationship row fires delete_orphaned_vector, which deletes the old vector.
create or replace function public.upsert_vector(p_embedding text, p_message_id uuid, p_todo_id uuid)
returns bigint
language plpgsql
as $$
declare
    _vector_id bigint;
begin
    delete from public.vector_relationship
    where (p_message_id is not null and message_id = p_message_id)
       or (p_todo_id is not null and todo_id = p_todo_id);

    insert into public.vector_table (embedding)
    values (p_embedding::extensions.vector)
    returning vector_id into _vector_id;

    insert into public.vector_relationship (vector_id, message_id, todo_id)
    values (_vector_id, p_message_id, p_todo_id);

    return _vector_id;
end;
$$;

-- only the edge function can call upsert_vector, not the public or authenticated roles
revoke execute on function public.upsert_vector(text, uuid, uuid) from public;
do $$
begin
    if exists (select 1 from pg_roles where rolname = 'anon') then
        revoke execute on function public.upsert_vector(text, uuid, uuid) from anon, authenticated;
    end if;
    if exists (select 1 from pg_roles where rolname = 'service_role') then
        grant execute on function public.upsert_vector(text, uuid, uuid) to service_role;
    end if;
end;
$$;

-- Sends the changed row to the embed-vectors edge function without blocking the insert/update.
create or replace function public.notify_embed_vectors()
returns trigger
language plpgsql
security definer
set search_path = ''
as $$
declare
    _url  text;
    _key  text;
    _body jsonb;
begin
    begin
        select decrypted_secret into _url from vault.decrypted_secrets where name = 'embed_vectors_function_url' limit 1;
        select decrypted_secret into _key from vault.decrypted_secrets where name = 'embed_vectors_function_key' limit 1;
    exception when undefined_table or invalid_schema_name then
        -- vault not available —> skip embedding
        return new;
    end;

    if _url is null or _key is null then
        return new;
    end if;

    _body := jsonb_build_object(
        'type',       TG_OP,
        'table',      TG_TABLE_NAME,
        'record',     row_to_json(new),
        'old_record', row_to_json(old)
    );

    begin
        perform net.http_post(
            url     := _url,
            headers := jsonb_build_object('Content-Type', 'application/json', 'Authorization', 'Bearer ' || _key),
            body    := _body
        );
    exception when others then
        -- never fail the insert/update because embedding couldn't be queued
        raise warning 'notify_embed_vectors failed for % on %, SQLSTATE %, error: %', TG_OP, TG_TABLE_NAME, SQLSTATE, SQLERRM;
    end;

    return new;
end;
$$;

-- Deletes need no trigger: the FK cascade removes the relationship row and delete_orphaned_vector removes the vector.
-- Updates only fire when the embedded text changes (not e.g. read_at or completed_at).
create trigger embed_chat_message_on_insert
after insert on public.chat_messages
for each row execute function public.notify_embed_vectors();

create trigger embed_chat_message_on_update
after update of message on public.chat_messages
for each row
when (old.message is distinct from new.message)
execute function public.notify_embed_vectors();

create trigger embed_todo_item_on_insert
after insert on public.todo_items
for each row execute function public.notify_embed_vectors();

create trigger embed_todo_item_on_update
after update of todo_description on public.todo_items
for each row
when (old.todo_description is distinct from new.todo_description)
execute function public.notify_embed_vectors();
