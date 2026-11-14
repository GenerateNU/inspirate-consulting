alter table public.activity
    add column created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
