-- create username index
CREATE INDEX username_idx ON users USING gin(name gin_trgm_ops);
-- create user_id index on counselor 
CREATE INDEX id_counselor_idx ON counselor (user_id);
-- create user_id index on student 
CREATE INDEX id_student_idx ON student (user_id);