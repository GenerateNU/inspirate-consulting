CREATE TYPE review_entry_type AS ENUM ('spend', 'refund', 'adjustment');
CREATE TYPE review_status AS ENUM ('open', 'completed', 'refunded');

CREATE TABLE IF NOT EXISTS essay_review_transaction (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subtotal     INT NOT NULL,
    student_id   UUID NOT NULL REFERENCES student(id),
    entry_type   review_entry_type NOT NULL,
    -- NULL for an adjustment. No FK yet: the essays table does not exist.
    essay_id     UUID,
    completed_at TIMESTAMPTZ,
    -- The row that reversed this charge. UNIQUE: a reversal serves one charge.
    refund       UUID UNIQUE REFERENCES essay_review_transaction(id),
    status       review_status GENERATED ALWAYS AS (
        CASE WHEN entry_type <> 'spend' THEN NULL
             WHEN refund IS NOT NULL THEN 'refunded'::review_status
             WHEN completed_at IS NOT NULL THEN 'completed'::review_status
             ELSE 'open'::review_status END
    ) STORED,
    actor_id     UUID REFERENCES users(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_essay_review_transaction_student_id ON essay_review_transaction (student_id);
