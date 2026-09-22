package userRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *UserRepository) FetchUser(ctx context.Context, user models.FetchUserInput) (*models.FetchUserOutput, error) {
	User := &models.FetchUserOutput{Body: &models.User{}}

	query, err := schema.ReadSQLBaseScript("get_user.sql", SqlUserFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	err = r.db.QueryRow(
		ctx,
		query,
		user.ID,
	).Scan(&User.Body.Name, &User.Body.SupabaseID, &User.Body.PfpKey)
	if err != nil {
		return nil, err
	}

	return User, nil
}
