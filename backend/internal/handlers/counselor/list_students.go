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

		if input.Year != "" && student.Year != input.Year {
			// skip students that don't match the year filter when it is set
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

// orders years from freshman to senior, since sorting the strings would be alphabetical
var yearRank = map[models.Year]int{
	models.Freshman:  1,
	models.Sophomore: 2,
	models.Junior:    3,
	models.Senior:    4,
}

// sorts by name, or by year then name
// ties fall back to student id
func sortStudents(students []models.CounselorStudent, sortBy string, sortOrder string) {

	slices.SortFunc(students, func(a, b models.CounselorStudent) int {

		byName := cmp.Compare(a.Name, b.Name)
		byYear := cmp.Compare(yearRank[a.Year], yearRank[b.Year])
		byID := cmp.Compare(a.StudentID.String(), b.StudentID.String())

		var result int

		switch sortBy {

		case "year":
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
