package counselorRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// fetches a counselor's profile by counselor id
func (r *CounselorRepository) GetCounselor(ctx context.Context, id uuid.UUID) (*models.Counselor, error) {
	counselor := &models.Counselor{}

	query, err := schema.ReadSQLBaseScript("get_counselor.sql", SqlCounselorFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, &err
	}

	err = r.db.QueryRow(ctx, query, id).Scan(&counselor.ID, &counselor.UserID, &counselor.CustomLink, &counselor.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("counselor", "id", id.String())
		}
		return nil, err
	}

	return counselor, nil
}
