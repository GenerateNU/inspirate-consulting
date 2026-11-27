package todoItemRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

// GetTodoItemStats counts a student's to-do items: total, incomplete, and due this week (Monday to Sunday, UTC).
func (r *TodoItemRepository) GetTodoItemStats(ctx context.Context, studentID string) (*models.TodoItemStats, error) {
	query, err := schema.ReadSQLBaseScript("get_todo_item_stats.sql", SqlTodoItemFiles)
	if err != nil {
		return nil, err
	}

	var stats models.TodoItemStats
	if err := r.db.QueryRow(ctx, query, studentID).Scan(
		&stats.TotalCount,
		&stats.IncompleteCount,
		&stats.DueThisWeekCount,
	); err != nil {
		return nil, err
	}

	return &stats, nil
}
