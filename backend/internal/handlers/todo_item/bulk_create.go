package todoitem

import (
	"context"
	"time"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (h *Handler) BulkCreateTodoItems(ctx context.Context, input *models.BulkCreateTodoItemsInput) (*models.BulkCreateTodoItemsResponse, error) {
	body := input.Body

	counselorID, err := auth.GetCounselorID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, errs.Unauthorized()
	}

	// Essays and extracurriculars belong to one student, so they can't be bulk-assigned
	if body.LinkType != nil &&
		(*body.LinkType == models.TodoLinkEssay || *body.LinkType == models.TodoLinkExtracurricular) {
		return nil, errs.BadRequest("bulk tasks can only link to a college or a module")
	}

	hasCollege := body.CollegeID != nil
	hasMedia := body.MediaID != nil
	if hasCollege && hasMedia {
		return nil, errs.BadRequest("a todo item can be linked to at most one of essay, media, college, or extracurricular")
	}
	if err := validateLink(body.LinkType, false, hasMedia, hasCollege, false); err != nil {
		return nil, err
	}

	if !body.All && len(body.StudentIDs) == 0 && len(body.GraduationYears) == 0 {
		return nil, errs.BadRequest("choose students, graduation years, or all")
	}

	assignmentID, items, err := h.TodoItemRepository.BulkCreateTodoItems(ctx, models.BulkCreateTodoItemsParams{
		CounselorID: counselorID,
		UserID:      userID,
		Title:       body.Title,
		Description: body.Description,
		Deadline:    body.Deadline,
		LinkType:    body.LinkType,
		CollegeID:   body.CollegeID,
		MediaID:     body.MediaID,
		StudentIDs:  body.StudentIDs,
		Years:       graduationYearsToStudentYears(body.GraduationYears, time.Now()),
		All:         body.All,
	})
	if err != nil {
		return nil, err
	}

	return &models.BulkCreateTodoItemsResponse{AssignmentID: assignmentID, Items: items}, nil
}

// graduationYearsToStudentYears converts graduation years to student years
// for the current school year. Students store a year, not a
// graduation year. The school year ends in the spring, so August onward counts
// toward the next calendar year.
func graduationYearsToStudentYears(gradYears []int, now time.Time) []string {
	schoolYearEnd := now.Year()
	if now.Month() >= time.August {
		schoolYearEnd++
	}
	byOffset := map[int]models.Year{
		0: models.Senior,
		1: models.Junior,
		2: models.Sophomore,
		3: models.Freshman,
	}
	var years []string
	for _, g := range gradYears {
		if y, ok := byOffset[g-schoolYearEnd]; ok {
			years = append(years, string(y))
		}
	}
	return years
}