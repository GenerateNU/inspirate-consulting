package extracurricularRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
)

// DeleteExtracurricular removes an extracurricular owned by studentID. Its
// achievements and notes are removed by the foreign keys' ON DELETE CASCADE.
// Returns a 404 HTTPError if no matching row exists for that student.
func (r *ExtracurricularRepository) DeleteExtracurricular(ctx context.Context, id int64, studentID string) error {
	deleteQuery, err := schema.ReadSQLBaseScript("delete_extracurricular.sql", SqlExtracurricularFiles)
	if err != nil {
		return err
	}

	result, err := r.db.Exec(ctx, deleteQuery, id, studentID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errs.NotFound("extracurricular", "id", id)
	}
	return nil
}
