package personalCollegeApplicationRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
)

func (r *PersonalCollegeApplicationRepository) DeletePersonalCollegeApplication(ctx context.Context, studentID string, applicationID int64) error {
	deleteQuery, err := schema.ReadSQLBaseScript("delete_personal_college_application.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return err
	}

	commandTag, err := r.db.Exec(ctx, deleteQuery, applicationID, studentID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errs.NotFound("personal college application", "id", applicationID)
	}

	return nil
}
