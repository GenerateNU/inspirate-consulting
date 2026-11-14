-- Required columns: a NULL parameter means "not sent", so COALESCE keeps the value.
-- Nullable columns: each takes a (sent boolean, value) pair. When sent is true
-- the value is written as-is, including NULL, which clears the column.
UPDATE public.extracurriculars SET
    name                    = COALESCE($1, name),
    status                  = COALESCE($2, status),
    type                    = COALESCE($3, type),
    start_date              = COALESCE($4, start_date),
    organization            = COALESCE($5, organization),
    ready_for_application   = COALESCE($6, ready_for_application),
    description             = CASE WHEN $7::boolean  THEN $8::text      ELSE description END,
    leadership_role         = CASE WHEN $9::boolean  THEN $10::text     ELSE leadership_role END,
    end_date                = CASE WHEN $11::boolean THEN $12::date     ELSE end_date END,
    hours_per_week          = CASE WHEN $13::boolean THEN $14::numeric  ELSE hours_per_week END,
    weeks_per_year          = CASE WHEN $15::boolean THEN $16::smallint ELSE weeks_per_year END,
    goals                   = CASE WHEN $17::boolean THEN $18::text     ELSE goals END,
    application_description = CASE WHEN $19::boolean THEN $20::text     ELSE application_description END,
    modified_at             = now()
WHERE id = $21
RETURNING
    id,
    created_at,
    modified_at,
    student_id,
    user_id,
    name,
    status,
    type,
    description,
    leadership_role,
    start_date,
    end_date,
    organization,
    hours_per_week,
    weeks_per_year,
    goals,
    ready_for_application,
    application_description;
