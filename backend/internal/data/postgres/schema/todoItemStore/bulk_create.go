package todoItemRepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// BulkCreateTodoItems assigns one task to many of a counselor's students in a single
// transaction. Every row shares one assignment_id, and module tasks also grant media access.
func (r *TodoItemRepository) BulkCreateTodoItems(ctx context.Context, p models.BulkCreateTodoItemsParams) (uuid.UUID, []models.TodoItem, error) {
	resolveQuery, err := schema.ReadSQLBaseScript("resolve_counselor_students.sql", SqlTodoItemFiles)
	if err != nil {
		return uuid.Nil, nil, err
	}
	insertQuery, err := schema.ReadSQLBaseScript("bulk_create_todo_items.sql", SqlTodoItemFiles)
	if err != nil {
		return uuid.Nil, nil, err
	}
	grantQuery, err := schema.ReadSQLBaseScript("grant_media_access_bulk.sql", SqlTodoItemFiles)
	if err != nil {
		return uuid.Nil, nil, err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, nil, err
	}
	defer tx.Rollback(ctx) // does nothing once committed

	// Empty slices instead of nil, so Postgres gets empty arrays rather than NULL
	studentIDs := p.StudentIDs
	if studentIDs == nil {
		studentIDs = []uuid.UUID{}
	}
	years := p.Years
	if years == nil {
		years = []string{}
	}

	// Resolve to only this counselor's students, deduplicated
	rows, err := tx.Query(ctx, resolveQuery, p.CounselorID, p.All, studentIDs, years)
	if err != nil {
		return uuid.Nil, nil, err
	}
	resolved, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return uuid.Nil, nil, err
	}
	if len(resolved) == 0 {
		return uuid.Nil, nil, errs.BadRequest("no matching students to assign")
	}

	// One row per student, all sharing assignment_id
	assignmentID := uuid.New()
	rows, err = tx.Query(ctx, insertQuery,
		p.UserID, p.Title, p.Description, p.Deadline, p.LinkType,
		p.CollegeID, p.MediaID, assignmentID, resolved)
	if err != nil {
		return uuid.Nil, nil, err
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.TodoItem])
	if err != nil {
		return uuid.Nil, nil, err
	}

	// Module tasks: grant access to the video, in the same transaction
	if p.MediaID != nil {
		if _, err := tx.Exec(ctx, grantQuery, *p.MediaID, resolved); err != nil {
			return uuid.Nil, nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, nil, err
	}
	return assignmentID, items, nil
}
