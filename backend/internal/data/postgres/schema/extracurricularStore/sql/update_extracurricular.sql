UPDATE public.extracurriculars SET
    name            = COALESCE($1, name),
    status          = COALESCE($2, status),
    type            = COALESCE($3, type),
    description     = COALESCE($4, description),
    leadership_role = COALESCE($5, leadership_role),
    start_date      = COALESCE($6, start_date),
    end_date        = COALESCE($7, end_date),
    organization    = COALESCE($8, organization),
    modified_at     = now()
WHERE id = $9
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
    organization;