package extracurricular

import (
	"context"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// ListExtracurriculars fetches every extracurricular for one student. A logged-in
// student always gets their own (from auth; requestedStudentID is ignored). A
// counselor must pass requestedStudentID for a student assigned to them.
func (h *Handler) ListExtracurriculars(ctx context.Context, requestedStudentID *string) ([]models.Extracurricular, error) {
	if studentID, err := auth.GetStudentID(ctx); err == nil {
		return h.ExtracurricularRepository.ListExtracurriculars(ctx, studentID.String())
	}

	if requestedStudentID == nil {
		return nil, errs.BadRequest("student_id is required")
	}
	studentID, err := auth.ParseStudentID(*requestedStudentID)
	if err != nil {
		return nil, err
	}
	if err := h.Authorizer.CanAccessStudent(ctx, studentID); err != nil {
		return nil, err
	}

	return h.ExtracurricularRepository.ListExtracurriculars(ctx, studentID.String())
}
