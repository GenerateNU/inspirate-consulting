package todoItemRepository

import (
	"context"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *TodoItemRepository) CreateTodoItem(ctx context.Context, item *models.CreateTodoItemRequestBody) (*models.TodoItem, error) {
	createdItem := &models.TodoItem{}

	insertQuery, err := schema.ReadSQLBaseScript("create_todo_item.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		insertQuery,
		item.StudentID,
		item.UserID,
		item.TodoDescription,
		item.Deadline,
		item.EssayID,
		item.MediaID,
		item.GlobalCollegeID,
	).Scan(
		&createdItem.ID,
		&createdItem.CreatedAt,
		&createdItem.UpdatedAt,
		&createdItem.StudentID,
		&createdItem.UserID,
		&createdItem.TodoDescription,
		&createdItem.CompletedAt,
		&createdItem.Deadline,
		&createdItem.EssayID,
		&createdItem.MediaID,
		&createdItem.GlobalCollegeID,
	)
	if err != nil {
		return nil, err
	}

	return createdItem, nil
}
