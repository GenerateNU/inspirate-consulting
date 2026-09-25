package userRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *UserRepository) CreateUser(ctx context.Context, user models.CreateUserInput, supabase_id uuid.UUID) (*models.CreateUserOutput, error) {
	createdUser := &models.CreateUserOutput{Body: &models.User{}}

	query, err := schema.ReadSQLBaseScript("create_user.sql", SqlUserFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		query,
		user.Body.Name,
		supabase_id,
		user.Body.PfpKey,
	).Scan(&createdUser.Body.ID, &createdUser.Body.Name, &createdUser.Body.SupabaseID, &createdUser.Body.PfpKey)
	if err != nil {
		fmt.Println("[db] QueryRow/Scan error:", err)
		return nil, err
	}

	return createdUser, nil
}
