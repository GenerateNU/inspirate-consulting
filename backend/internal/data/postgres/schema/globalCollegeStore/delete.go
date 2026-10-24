package globalCollegeRepository

import (
	"context"
	"errors"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *GlobalCollegeRepository) DeleteGlobalCollege(ctx context.Context, id int64) error {
	deleteQuery, err := schema.ReadSQLBaseScript("delete_global_college.sql", SqlGlobalCollegeFiles)
	if err != nil {
		return err
	}

	commandTag, err := r.db.Exec(ctx, deleteQuery, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return errs.BadRequest("cannot delete: this college is still in use")
		}
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errs.NotFound("global college", "id", id)
	}

	return nil
}
