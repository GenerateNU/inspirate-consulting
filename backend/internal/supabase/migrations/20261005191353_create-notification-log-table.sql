CREATE TYPE NOTIFICATION_TYPE AS ENUM ('upcoming', 'past_due');

create table IF NOT EXISTS notification_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    task_id UUID REFERENCES todo_items(id),
    notification_type NOTIFICATION_TYPE NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);