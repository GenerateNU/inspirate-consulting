
CREATE TYPE YEAR AS ENUM ('freshman', 'sophomore', 'junior', 'senior');

CREATE TABLE IF NOT EXISTS user {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    supabase_id UUID gen_random_uuid(),
    pfp_key TEXT
};

CREATE TABLE IF NOT EXISTS student {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES user(id),
    year YEAR NOT NULL,
    organization TEXT,
    review_balance INT NOT NULL,
    counselor_id UUID NOT NULL REFERENCES counselor(id)
};

CREATE TABLE IF NOT EXISTS counselor {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES user(id)
};