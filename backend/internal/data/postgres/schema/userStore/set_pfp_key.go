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

// sets the user's pfp_key (nil clears it) and returns the key it replaced
func (r *UserRepository) SetPfpKey(ctx context.Context, id uuid.UUID, pfpKey *string) (*models.User, *string, error) {
	updatedUser := &models.User{}
	var previousKey *string

	query, err := schema.ReadSQLBaseScript("set_pfp_key.sql", SqlUserFiles)
	if err != nil {
		return nil, nil, errs.InternalServerError("Failed to read base query: ", err.Error())
	}

	err = r.db.QueryRow(ctx, query, id, pfpKey).Scan(
		&updatedUser.ID, &updatedUser.Name, &updatedUser.SupabaseID, &updatedUser.PfpKey,
		&updatedUser.ResetTime, &updatedUser.Email, &updatedUser.ExpiresAt, &previousKey,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, errs.NotFound("user", "id", id.String())
		}
		return nil, nil, err
	}

	return updatedUser, previousKey, nil
}
