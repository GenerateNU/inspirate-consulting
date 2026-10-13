package userRepository

import (
	"context"
	"errors"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *UserRepository) FetchUserBySupabaseID(ctx context.Context, input models.FetchUserBySupabaseIDInput) (*models.FetchUserOutput, error) {
	User := &models.FetchUserOutput{Body: &models.User{}}

	query, err := schema.ReadSQLBaseScript("get_user_by_supabase_id.sql", SqlUserFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	err = r.db.QueryRow(
		ctx,
		query,
		input.SupabaseID,
	).Scan(&User.Body.ID, &User.Body.Name, &User.Body.SupabaseID, &User.Body.PfpKey, &User.Body.ResetTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("user", "supabase_id", input.SupabaseID.String())
		}
		return nil, err
	}

	return User, nil
}
