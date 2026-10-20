package studentstore

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
)

func (r *StudentRepository) AssignCounselor(ctx context.Context, studentUserID uuid.UUID, counselorID uuid.UUID) (*models.Student, error) {
	query, err := schema.ReadSQLBaseScript("assign_student_to_counselor.sql", SqlStudentFiles)
	if err != nil {
		return nil, err
	}

	student := &models.Student{}
	err = r.db.QueryRow(ctx, query, studentUserID, counselorID).Scan(
		&student.ID, &student.UserID, &student.Year, &student.Organization, &student.GPA, &student.ReviewBalance, &student.CounselorID,
	)
	if err != nil {
		return nil, err
	}

	return student, nil
}
