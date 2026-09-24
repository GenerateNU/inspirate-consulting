package todoItemRepository

import "github.com/jackc/pgx/v5/pgxpool"

type TodoItemRepository struct {
	db *pgxpool.Pool
}

func NewTodoItemRepository(db *pgxpool.Pool) *TodoItemRepository {
	return &TodoItemRepository{db: db}
}
