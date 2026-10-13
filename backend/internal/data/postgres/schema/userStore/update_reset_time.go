package userRepository

import (
	"context"
	"time"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"

	"github.com/google/uuid"
)

func (r *UserRepository) UpdateResetTime(ctx context.Context, supabaseID uuid.UUID, resetTime *time.Time) error {
	query, err := schema.ReadSQLBaseScript("update_reset_time.sql", SqlUserFiles)
	if err != nil {
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return &err
	}

	_, err = r.db.Exec(ctx, query, supabaseID, resetTime)
	return err
}
