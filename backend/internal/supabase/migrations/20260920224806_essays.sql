CREATE TYPE status AS ENUM ('Submitted', 'Draft', 'Review', 'Archived');

CREATE TABLE essays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    /* this is a fk table does not exist yet though so its a place holder*/
    student_id UUID NOT NULL DEFAULT gen_random_uuid() /*REFERENCES students(id)*/,
    type TEXT NOT NULL,
    college_id BIGINT REFERENCES global_colleges(id),
    link_to_content TEXT NOT NULL,
    status status NOT NULL DEFAULT 'Draft'
);