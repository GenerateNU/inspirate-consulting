

/* new essay group table */
CREATE TABLE essay_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
/* To Do make this refrence the students table */
    student_id UUID NOT NULL
);

/* add new group_id column to essay */  

ALTER TABLE essays 
ADD COLUMN essay_group_id UUID REFERENCES essay_groups(id)