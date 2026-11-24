alter table public.media
    add column youtube_video_id text,
    alter column length_in_mins drop not null,
    alter column s3_key drop not null,
    add constraint media_single_source check (num_nonnulls(s3_key, youtube_video_id) = 1);