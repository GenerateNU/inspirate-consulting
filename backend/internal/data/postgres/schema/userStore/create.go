package userRepository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *UserRepository) CreateUser(ctx context.Context, user models.CreateUserInput, supabase_id uuid.UUID) (*models.CreateUserOutput, error) {
	createdUser := &models.User{}
	oneSecondAgo := time.Now().Add(-1 * time.Second)

	query, err := schema.ReadSQLBaseScript("create_user.sql", SqlUserFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	err = r.db.QueryRow(
		ctx,
		query,
		user.Body.Name,
		supabase_id,
		user.Body.PfpKey,
		oneSecondAgo,
	).Scan(&createdUser.ID, &createdUser.Name, &createdUser.SupabaseID, &createdUser.PfpKey, &createdUser.ResetTime)
	if err != nil {
		return nil, err
	}

	return &models.CreateUserOutput{Body: &models.CreateUserBody{
		User: createdUser,
	}}, nil
}
