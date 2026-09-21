package userRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

<<<<<<< HEAD
func (r *UserRepository) FetchUser(ctx context.Context, user models.FetchUserInput) (*models.FetchUserOutput, error) {
	User := &models.FetchUserOutput{Body: &models.User{}}
=======
func (r *UserRepository) FetcheUser(ctx context.Context, user models.FetchUserInput) (*models.FetchUserOutput, error) {
	User := &models.FetchUserOutput{}
>>>>>>> 4d96216 (added get user flow)

	query, err := schema.ReadSQLBaseScript("get_user.sql", SqlUserFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	err = r.db.QueryRow(
		ctx,
		query,
<<<<<<< HEAD
		user.ID,
	).Scan(&User.Body.Name, &User.Body.SupabaseID, &User.Body.PfpKey)
=======
		User,
	).Scan(&User.Body.ID, &User.Body.Name, &User.Body.SupabaseID, &User.Body.PfpKey)
>>>>>>> 4d96216 (added get user flow)
	if err != nil {
		return nil, err
	}

	return User, nil
}
