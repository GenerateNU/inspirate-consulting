package userRepository

import (
	"context"
	"errors"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *UserRepository) FetchSessionUserBySupabaseID(ctx context.Context, supabaseID uuid.UUID) (*models.SessionUser, error) {
	user := &models.SessionUser{}

	query, err := schema.ReadSQLBaseScript("get_session_user_by_supabase_id.sql", SqlUserFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	err = r.db.QueryRow(ctx, query, supabaseID).Scan(
		&user.ID, &user.Name, &user.SupabaseID, &user.PfpKey, &user.ResetTime, &user.StudentID, &user.CounselorID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("user", "supabase_id", supabaseID.String())
		}
		return nil, err
	}

	return user, nil
}
