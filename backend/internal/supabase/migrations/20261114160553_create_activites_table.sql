CREATE TYPE action AS ENUM (
    'added_school',
    'removed_school',
    'essay_status_changed',
    'added_extracurricular',
    'removed_extracurricular',
    'task_past_due'
);

CREATE TABLE public.activity (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES student(id),
    action action NOT NULL,
    essay_id UUID REFERENCES essays(id),
    essay_status status,
    global_college_id BIGINT REFERENCES global_colleges(id),
    extracurricular_id BIGINT REFERENCES extracurriculars(id),
    todo_item_id UUID REFERENCES todo_items(id)
);
