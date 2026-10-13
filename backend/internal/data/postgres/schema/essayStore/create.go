package essayRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayRepository) CreateEssay(ctx context.Context, essay models.Essays) error {
	insertQuery, err := schema.ReadSQLBaseScript("create_essay.sql", SqlEssayFiles)
	if err != nil {
		return errs.InternalServerError("Failed to read base query: ", err.Error())
	}

	_, err = r.db.Exec(ctx, insertQuery, essay.StudentID, essay.Type, essay.CollegeID, essay.LinkToContent, essay.EssayGroupID)

	if err != nil {
		var pgErr *pgconn.PgError

		// essay_group_id references essay_groups(id), so an unknown group is a client error
		if errors.As(err, &pgErr) && pgErr.Code == "23503" && essay.EssayGroupID != nil {
			return errs.NotFound("essay group", "id", essay.EssayGroupID.String())
		}

		return err
	}

	return nil
}
