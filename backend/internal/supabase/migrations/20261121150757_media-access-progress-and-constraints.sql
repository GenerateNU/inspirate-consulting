alter table public.media_access
    add column started_at timestamptz,
    add column completed_at timestamptz,
    unique (student_id, media_id),
    drop constraint media_access_media_id_fkey,
    add constraint media_access_media_id_fkey
        foreign key (media_id) references public.media(id) on delete cascade;