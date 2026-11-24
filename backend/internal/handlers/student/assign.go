package student

import (
	"context"

	"inspirate-consulting/internal/models"
)

// AssignCounselor sets or clears (nil CounselorID) a student's counselor.
// If the student already has a counselor, only that counselor may change it.
// An unassigned student can be picked up by any counselor.
func (h *Handler) AssignCounselor(ctx context.Context, input *models.AssignCounselorInput) (*models.Student, error) {
	studentID, err := h.StudentRepository.GetStudentIDByUserID(ctx, input.Body.UserID)
	if err != nil {
		return nil, err
	}

	student, err := h.StudentRepository.GetStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	if student.CounselorID != nil {
		if err := h.Authorizer.CanAccessStudent(ctx, studentID); err != nil {
			return nil, err
		}
	}

	return h.StudentRepository.AssignCounselor(ctx, input.Body.UserID, input.Body.CounselorID)
}
