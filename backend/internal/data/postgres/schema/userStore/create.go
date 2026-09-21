package userRepository

import (
	"context"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *UserRepository) CreateUser(ctx context.Context, greeting models.CreateUserInput) (*models.CreateUserOutput, error) {
	createdUser := &models.CreateUserOutput{}

	query, err := schema.ReadSQLBaseScript("create_user.sql", SqlEventFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	err = r.db.QueryRow(
		ctx,
		query,
		greeting.Body.Name,
		createdUser,
	).Scan(&createdUser.Body.ID, &createdUser.Body.Name, &createdUser.Body.SupabaseID, &createdUser.Body.PfpKey)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
