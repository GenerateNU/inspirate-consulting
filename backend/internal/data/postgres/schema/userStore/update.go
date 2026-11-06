package userRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// updates a user's name and pfp_key, leaving any nil field unchanged
func (r *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, input models.UpdateUserBody) (*models.User, error) {
	updatedUser := &models.User{}

	query, err := schema.ReadSQLBaseScript("update_user.sql", SqlUserFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	err = r.db.QueryRow(
		ctx,
		query,
		id,
		input.Name,
		input.PfpKey,
	).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.SupabaseID, &updatedUser.PfpKey, &updatedUser.ResetTime, &updatedUser.Email, &updatedUser.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("user", "id", id.String())
		}
		return nil, err
	}

	return updatedUser, nil
}
