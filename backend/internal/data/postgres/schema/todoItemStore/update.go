package todoItemRepository

import (
	"context"
	"time"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *TodoItemRepository) UpdateTodoItemCompletedAt(ctx context.Context, id string, completedAt *time.Time) (*models.TodoItem, error) {
	updatedItem := &models.TodoItem{}

	updateQuery, err := schema.ReadSQLBaseScript("update_todo_item_completed_at.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		updateQuery,
		completedAt,
		id,
	).Scan(
		&updatedItem.ID,
		&updatedItem.CreatedAt,
		&updatedItem.UpdatedAt,
		&updatedItem.StudentID,
		&updatedItem.UserID,
		&updatedItem.TodoDescription,
		&updatedItem.CompletedAt,
		&updatedItem.Deadline,
		&updatedItem.EssayID,
		&updatedItem.MediaID,
		&updatedItem.GlobalCollegeID,
	)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

func (r *TodoItemRepository) UpdateTodoItemDescription(ctx context.Context, id string, newDescription string) (*models.TodoItem, error) {
	updatedItem := &models.TodoItem{}

	updateQuery, err := schema.ReadSQLBaseScript("update_todo_item_description.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		updateQuery,
		id,
		newDescription,
	).Scan(
		&updatedItem.ID,
		&updatedItem.CreatedAt,
		&updatedItem.UpdatedAt,
		&updatedItem.StudentID,
		&updatedItem.UserID,
		&updatedItem.TodoDescription,
		&updatedItem.CompletedAt,
		&updatedItem.Deadline,
		&updatedItem.EssayID,
		&updatedItem.MediaID,
		&updatedItem.GlobalCollegeID,
	)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

func (r *TodoItemRepository) UpdateTodoItemDeadline(ctx context.Context, id string, newDeadline *time.Time) (*models.TodoItem, error) {
	updatedItem := &models.TodoItem{}

	updateQuery, err := schema.ReadSQLBaseScript("update_todo_item_deadline.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		updateQuery,
		id,
		newDeadline,
	).Scan(
		&updatedItem.ID,
		&updatedItem.CreatedAt,
		&updatedItem.UpdatedAt,
		&updatedItem.StudentID,
		&updatedItem.UserID,
		&updatedItem.TodoDescription,
		&updatedItem.CompletedAt,
		&updatedItem.Deadline,
		&updatedItem.EssayID,
		&updatedItem.MediaID,
		&updatedItem.GlobalCollegeID,
	)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}
