package userRepository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// creates the users row and its student/counselor row in one transaction
func (r *UserRepository) CreateUser(ctx context.Context, user models.CreateUserInput, supabase_id uuid.UUID) (*models.CreateUserOutput, error) {
	createdUser := &models.User{}
	oneSecondAgo := time.Now().Add(-1 * time.Second)

	query, err := schema.ReadSQLBaseScript("create_user.sql", SqlUserFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	// Run as a transaction so a user is never left without their student/counselor row
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	err = tx.QueryRow(
		ctx,
		query,
		user.Body.Name,
		supabase_id,
		user.Body.PfpKey,
		oneSecondAgo,
		user.Body.Email,
		user.Body.ExpiresAt,
	).Scan(&createdUser.ID, &createdUser.Name, &createdUser.SupabaseID, &createdUser.PfpKey, &createdUser.ResetTime, &createdUser.Email, &createdUser.ExpiresAt)

	if err != nil {
		return nil, err
	}

	switch user.Body.Role {
	case models.StudentRole:
		studentQuery, err := schema.ReadSQLBaseScript("create_student.sql", SqlUserFiles)
		if err != nil {
			err := errs.InternalServerError("Failed to read base query: ", err.Error())
			return nil, &err
		}
		_, err = tx.Exec(
			ctx,
			studentQuery,
			createdUser.ID,
			user.Body.Year,
			user.Body.Organization,
			user.Body.GPA,
			user.Body.CounselorID,
			user.Body.Interest,
		)
		if err != nil {
			return nil, err
		}
	case models.CounselorRole:
		counselorQuery, err := schema.ReadSQLBaseScript("create_counselor.sql", SqlUserFiles)
		if err != nil {
			err := errs.InternalServerError("Failed to read base query: ", err.Error())
			return nil, &err
		}
		_, err = tx.Exec(ctx, counselorQuery, createdUser.ID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errs.BadRequest("role must be student or counselor")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.CreateUserOutput{Body: &models.CreateUserBody{
		User: createdUser,
	}}, nil
}
