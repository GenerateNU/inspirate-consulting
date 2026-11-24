package todoItemRepository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// UpdateTodoAssignment edits every task in an assignment and adds or removes students, in one transaction. Removing a student keeps their task if they already completed it.
func (r *TodoItemRepository) UpdateTodoAssignment(ctx context.Context, p models.UpdateTodoAssignmentParams) ([]models.TodoItem, error) {
	read := func(name string) (string, error) {
		return schema.ReadSQLBaseScript(name, SqlTodoItemFiles)
	}
	countQuery, err := read("count_assignment_rows.sql")
	if err != nil {
		return nil, err
	}
	updateQuery, err := read("update_assignment_fields.sql")
	if err != nil {
		return nil, err
	}
	addQuery, err := read("add_assignment_students.sql")
	if err != nil {
		return nil, err
	}
	removeQuery, err := read("remove_assignment_students.sql")
	if err != nil {
		return nil, err
	}
	studentIDsQuery, err := read("assignment_student_ids.sql")
	if err != nil {
		return nil, err
	}
	grantQuery, err := read("grant_media_access_bulk.sql")
	if err != nil {
		return nil, err
	}
	getQuery, err := read("get_assignment_items.sql")
	if err != nil {
		return nil, err
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// The assignment must exist and belong to this counselor
	var count int
	if err := tx.QueryRow(ctx, countQuery, p.AssignmentID, p.UserID).Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errs.NotFound("todo assignment", "id", p.AssignmentID)
	}

	// Edit the shared fields on every row
	if _, err := tx.Exec(ctx, updateQuery,
		p.AssignmentID, p.Title, p.Description, p.Deadline,
		p.LinkType, p.CollegeID, p.MediaID); err != nil {
		return nil, err
	}

	// Add students (after the edit, so new rows copy the updated fields)
	if len(p.AddStudentIDs) > 0 {
		if _, err := tx.Exec(ctx, addQuery, p.AssignmentID, p.CounselorID, p.AddStudentIDs); err != nil {
			return nil, err
		}
	}

	// Remove students, keeping completed tasks
	if len(p.RemoveStudentIDs) > 0 {
		if _, err := tx.Exec(ctx, removeQuery, p.AssignmentID, p.RemoveStudentIDs); err != nil {
			return nil, err
		}
	}

	// Read back the assignment's rows
	rows, err := tx.Query(ctx, getQuery, p.AssignmentID)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.TodoItem])
	if err != nil {
		return nil, err
	}

	// Module tasks: make sure everyone on the assignment has access to the video
	if len(items) > 0 && items[0].MediaID != nil {
		rows, err := tx.Query(ctx, studentIDsQuery, p.AssignmentID)
		if err != nil {
			return nil, err
		}
		studentIDs, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
		if err != nil {
			return nil, err
		}
		if _, err := tx.Exec(ctx, grantQuery, *items[0].MediaID, studentIDs); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return items, nil
}
