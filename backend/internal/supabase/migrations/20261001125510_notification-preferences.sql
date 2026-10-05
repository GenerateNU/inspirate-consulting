
create table public.notification_preferences (
    user_id uuid primary key references public.users(id),
    email_enabled boolean not null default true,
    -- weekly reminders of outstanding assignments, etc.
    weekly_summary_enabled boolean not null default true,

    due_date_notifications_enabled boolean not null default true,
    -- send a notification to the user that the assignment is due this many days before the due date
    days_before_due integer not null check (days_before_due >= 0) default 1,

    notify_past_due boolean not null default true
);