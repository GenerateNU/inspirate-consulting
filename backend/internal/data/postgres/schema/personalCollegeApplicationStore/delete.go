package personalCollegeApplicationRepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
)

// removes a school from the student's list, in one transaction:
//   - deletes the application
//   - deletes the essays the student started for that college (essays with review history are archived
//     instead, since the credit ledger references them; essays with no college are never touched)
//   - closes the gap in the student's ranking
//
// essays are left alone if the student still has another application to the same college.
func (r *PersonalCollegeApplicationRepository) DeletePersonalCollegeApplication(ctx context.Context, studentID uuid.UUID, applicationID int64) error {
	deleteQuery, err := schema.ReadSQLBaseScript("delete_personal_college_application.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return err
	}
	archiveEssaysQuery, err := schema.ReadSQLBaseScript("archive_college_essays_with_reviews.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return err
	}
	deleteEssaysQuery, err := schema.ReadSQLBaseScript("delete_college_essays.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return err
	}
	closeRankGapQuery, err := schema.ReadSQLBaseScript("close_personal_college_application_rank_gap.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var collegeID int64
	var rank *int
	if err := tx.QueryRow(ctx, deleteQuery, applicationID, studentID).Scan(&collegeID, &rank); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.NotFound("personal college application", "id", applicationID)
		}
		return err
	}

	if _, err := tx.Exec(ctx, archiveEssaysQuery, studentID, collegeID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, deleteEssaysQuery, studentID, collegeID); err != nil {
		return err
	}

	if rank != nil {
		if _, err := tx.Exec(ctx, closeRankGapQuery, studentID, *rank); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
