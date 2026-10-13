package essayGroupRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *EssayGroupRepository) CreateEssayGroup(ctx context.Context, group models.CreateEssayGroupBody) (*models.EssayGroups, error) {
	createdEssayGroup := &models.EssayGroups{}

	insertQuery, err := schema.ReadSQLBaseScript("create_essay_group.sql", SqlEssayGroupFiles)
	if err != nil {
		return nil, errs.InternalServerError("Failed to read base query: ", err.Error())
	}

	err = r.db.QueryRow(
		ctx,
		insertQuery,
		group.Name,
		group.Description,
		group.StudentID,
	).Scan(
		&createdEssayGroup.ID,
		&createdEssayGroup.Name,
		&createdEssayGroup.Description,
		&createdEssayGroup.StudentID,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		// name is UNIQUE, so a duplicate is a conflict rather than a server error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errs.Conflict("essay group", "name", group.Name)
		}

		return nil, err
	}

	return createdEssayGroup, nil
}
