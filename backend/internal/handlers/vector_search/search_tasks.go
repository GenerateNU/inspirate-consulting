package vectorsearch

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// SearchTasks finds todo items closest in meaning to the query.
// Students search their own todos; counselors search the todos of all their assigned students.
func (h *Handler) SearchTasks(ctx context.Context, query string) ([]models.TodoItemSearchResult, error) {
	if h.EmbeddingClient == nil {
		return nil, errs.NewHTTPError(http.StatusServiceUnavailable, errors.New("similarity search is not configured"))
	}
	if strings.TrimSpace(query) == "" {
		return nil, errs.BadRequest("query cannot be blank")
	}

	studentIDs, err := h.visibleStudentIDs(ctx)
	if err != nil {
		return nil, err
	}
	if len(studentIDs) == 0 {
		return []models.TodoItemSearchResult{}, nil
	}

	queryEmbedding, err := h.EmbeddingClient.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	return h.VectorRepository.SearchTodoItems(ctx, queryEmbedding, studentIDs, searchLimit)
}

// visibleStudentIDs returns the students whose todos the current user may search.
func (h *Handler) visibleStudentIDs(ctx context.Context) ([]uuid.UUID, error) {
	switch models.Role(auth.GetRole(ctx)) {
	case models.StudentRole:
		studentID, err := auth.GetStudentID(ctx)
		if err != nil {
			return nil, errs.Unauthorized()
		}
		return []uuid.UUID{studentID}, nil

	case models.CounselorRole:
		counselorID, err := auth.GetCounselorID(ctx)
		if err != nil {
			return nil, errs.Unauthorized()
		}
		students, err := h.CounselorRepository.ListStudents(ctx, counselorID)
		if err != nil {
			return nil, err
		}
		ids := make([]uuid.UUID, len(students))
		for i, student := range students {
			ids[i] = student.StudentID
		}
		return ids, nil

	default:
		return nil, errs.Forbidden("insufficient role")
	}
}
