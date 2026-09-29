create table public.chat_messages (
    id uuid primary key default gen_random_uuid(),
    sender_id uuid not null references public.users(id),
    recipient_id uuid not null references public.users(id),
    message text not null check (length(trim(message)) > 0),
    created_at timestamptz not null default now(),
    edited_at timestamptz,
    read_at timestamptz,

    -- These two columns are used to create a unique index for each conversation between two users.
    user_low uuid generated always as (least(sender_id, recipient_id)) stored,
    user_high uuid generated always as (greatest(sender_id, recipient_id)) stored,

    check (sender_id <> recipient_id)
);

create index idx_chat_messages_conversation
on public.chat_messages (user_low, user_high, created_at desc, id desc);

create index idx_chat_messages_user_high
on public.chat_messages (user_high, created_at desc);

create index idx_chat_messages_unread
on public.chat_messages (recipient_id, sender_id)
where read_at is null;
