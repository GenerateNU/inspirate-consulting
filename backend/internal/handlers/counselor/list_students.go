package counselor

import (
	"cmp"
	"context"
	"slices"

	"github.com/google/uuid"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/models"
)

// lists the authenticated counselor's students, flagged for falling behind and pending reviews,
// then filters and sorts them according to the input
func (h *Handler) ListStudents(ctx context.Context, input models.ListCounselorStudentsInput) ([]models.CounselorStudent, error) {
	counselorID, err := auth.GetCounselorID(ctx)
	if err != nil {
		return nil, err
	}

	allStudents, err := h.CounselorRepository.ListStudents(ctx, counselorID)
	if err != nil {
		return nil, err
	}

	reviews, err := h.CounselorRepository.ListPendingReviews(ctx, counselorID)
	if err != nil {
		return nil, err
	}

	// map student IDs to whether they have a pending review
	hasPendingReview := map[uuid.UUID]bool{}
	for _, review := range reviews {
		hasPendingReview[review.StudentID] = true
	}

	fallingBehindIDs, err := h.CounselorRepository.ListFallingBehindStudentIDs(ctx, counselorID)
	if err != nil {
		return nil, err
	}

	students := make([]models.CounselorStudent, 0, len(allStudents))

	for _, student := range allStudents {

		if input.GraduationYear != 0 && student.GraduationYear != input.GraduationYear {
			// skip students that don't match the graduation year filter when it is set
			continue
		}

		if input.HasPendingReview && !hasPendingReview[student.StudentID] {
			// skip students that don't have a pending review when it is set
			continue
		}

		students = append(students, models.CounselorStudent{
			CounselorStudentRow: student,
			IsFallingBehind:     slices.Contains(fallingBehindIDs, student.StudentID),
			HasPendingReview:    hasPendingReview[student.StudentID],
		})
	}

	sortStudents(students, input.SortBy, input.SortOrder)

	return students, nil
}

// sorts by name, or by graduation year then name
// ties fall back to student id
func sortStudents(students []models.CounselorStudent, sortBy string, sortOrder string) {

	slices.SortFunc(students, func(a, b models.CounselorStudent) int {

		byName := cmp.Compare(a.Name, b.Name)
		byYear := cmp.Compare(a.GraduationYear, b.GraduationYear)
		byID := cmp.Compare(a.StudentID.String(), b.StudentID.String())

		var result int

		switch sortBy {

		case "graduation_year":
			result = cmp.Or(byYear, byName, byID)

		case "name":
			result = cmp.Or(byName, byID)
		}

		if sortOrder == "desc" {
			return -result
		}

		return result
	})
}
